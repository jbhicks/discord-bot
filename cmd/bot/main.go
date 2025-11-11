package main

import (
	"log/slog"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/josh/discord-bot/internal/db"
	"github.com/josh/discord-bot/internal/llm"
	"github.com/josh/discord-bot/internal/sentiment"
	"github.com/josh/discord-bot/internal/stocknews"
	"github.com/josh/discord-bot/pkg/commands"
)

var commandMap = make(map[string]commands.Command)

func registerCommands() {
	ping := &commands.PingCommand{}
	commandMap[ping.Name()] = ping
	play := &commands.PlayCommand{}
	commandMap[play.Name()] = play
	stop := &commands.StopCommand{}
	commandMap[stop.Name()] = stop
	queue := &commands.QueueCommand{}
	commandMap[queue.Name()] = queue
	skip := &commands.SkipCommand{}
	commandMap[skip.Name()] = skip
	loop := &commands.LoopCommand{}
	commandMap[loop.Name()] = loop
	search := &commands.SearchCommand{}
	commandMap[search.Name()] = search
	playlist := &commands.PlaylistCommand{}
	commandMap[playlist.Name()] = playlist
	ai := &commands.AICommand{}
	commandMap[ai.Name()] = ai
	help := &commands.HelpCommand{}
	commandMap[help.Name()] = help

	imageGenURL := os.Getenv("IMAGE_GEN_URL")
	if imageGenURL == "" {
		imageGenURL = "http://localhost:7860"
	}
	imagine := commands.NewImagineCommand(imageGenURL)
	commandMap[imagine.Name()] = imagine

	llmURL := os.Getenv("LLM_URL")
	if llmURL == "" {
		llmURL = "http://localhost:8081"
	}
	pdf := commands.NewPDFCommand(llmURL, imageGenURL)
	commandMap[pdf.Name()] = pdf

	marketauxAPIKey := os.Getenv("MARKETAUX_API_KEY")
	alphaVantageAPIKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	marketaux := stocknews.NewMarketAuxClient(marketauxAPIKey)
	alphavantage := stocknews.NewAlphaVantageClient(alphaVantageAPIKey)
	newsClient := stocknews.NewFallbackClient(marketaux, alphavantage)
	llmClient := llm.NewClient(llmURL)
	sentimentClient := sentiment.NewAggregator()
	stock := commands.NewStockCommand(newsClient, llmClient, sentimentClient)
	commandMap[stock.Name()] = stock
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}

	registerCommands()

	err = db.InitDB()
	if err != nil {
		slog.Error("Error initializing DB", "error", err)
		return
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		slog.Error("TOKEN not found")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Error("Error creating Discord session", "error", err)
		return
	}

	dg.AddHandler(ready)
	dg.AddHandler(interactionCreate)

	// Add intents for guilds and voice states
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates

	err = dg.Open()
	if err != nil {
		slog.Error("Error opening connection", "error", err)
		return
	}
	defer dg.Close()

	slog.Info("Bot is now running. Press CTRL-C to exit.")
	<-make(chan struct{})
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	slog.Info("Bot is ready", "user", s.State.User.Username)

	guildID := os.Getenv("GUILD_ID")
	if guildID == "" {
		guildID = "414275056265330689"
	}

	// // Delete all existing commands
	// commands, err := s.ApplicationCommands(s.State.User.ID, guildID)
	// if err != nil {
	// 	slog.Error("Cannot fetch commands", "error", err)
	// } else {
	// 	for _, cmd := range commands {
	// 		err := s.ApplicationCommandDelete(s.State.User.ID, guildID, cmd.ID)
	// 		if err != nil {
	// 			slog.Error("Cannot delete command", "name", cmd.Name, "error", err)
	// 		} else {
	// 			slog.Info("Deleted old command", "name", cmd.Name)
	// 		}
	// 	}
	// }

	// Register new commands
	for _, cmd := range commandMap {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd.Data())
		if err != nil {
			slog.Error("Cannot create command", "name", cmd.Name(), "error", err)
		} else {
			slog.Info("Registered command", "name", cmd.Name())
		}
	}
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	cmd, ok := commandMap[data.Name]
	if !ok {
		slog.Error("Unknown command", "name", data.Name)
		return
	}

	slog.Info("Executing command", "name", cmd.Name(), "user", i.Member.User.Username)

	err := cmd.Execute(s, i)
	if err != nil {
		slog.Error("Error executing command", "name", cmd.Name(), "error", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There was an error while executing this command!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	} else {
		slog.Info("Command executed successfully", "name", cmd.Name())
	}
}
