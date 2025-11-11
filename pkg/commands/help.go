package commands

import (
	"github.com/bwmarrin/discordgo"
)

type HelpCommand struct{}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Description() string {
	return "Show all available commands"
}

func (c *HelpCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "Show all available commands",
	}
}

func (c *HelpCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Since commandMap is in main, we can't access it here.
	// For simplicity, hardcode the list.
	content := "**Available Commands:**\n" +
		"- `/ping`: Replies with Pong!\n" +
		"- `/play <url>`: Play a song from MP3 URL\n" +
		"- `/stop`: Stop playing and leave voice channel\n" +
		"- `/queue add <url>`: Add song to queue\n" +
		"- `/queue remove <index>`: Remove song from queue\n" +
		"- `/queue view`: View current queue\n" +
		"- `/queue shuffle`: Shuffle the queue\n" +
		"- `/skip`: Skip current song\n" +
		"- `/loop`: Toggle loop mode\n" +
		"- `/search <query>`: Search YouTube for songs\n" +
		"- `/playlist create <name>`: Create a playlist\n" +
		"- `/playlist add <name> <url>`: Add song to playlist\n" +
		"- `/playlist play <name>`: Play a playlist\n" +
		"- `/playlist list`: List your playlists\n" +
		"- `/ai <prompt>`: Generate AI content\n" +
		"- `/imagine <prompt>`: Generate images with Stable Diffusion\n" +
		"- `/pdf`: Generate PDF documents with AI\n" +
		"  • Types: Document/Report, Presentation/Slides, Spreadsheet/Table\n" +
		"  • Automatically includes AI-generated images\n" +
		"- `/help`: Show this help"

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}
