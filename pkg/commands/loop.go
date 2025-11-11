package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/josh/discord-bot/internal/voice"
)

type LoopCommand struct{}

func (c *LoopCommand) Name() string {
	return "loop"
}

func (c *LoopCommand) Description() string {
	return "Toggle loop mode for the queue"
}

func (c *LoopCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "loop",
		Description: "Toggle loop mode for the queue",
	}
}

func (c *LoopCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	current := voice.IsLooping(i.GuildID)
	newState := !current
	voice.SetLoop(i.GuildID, newState)
	status := "disabled"
	if newState {
		status = "enabled"
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Loop " + status,
		},
	})
}
