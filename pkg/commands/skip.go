package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/josh/discord-bot/internal/voice"
)

type SkipCommand struct{}

func (c *SkipCommand) Name() string {
	return "skip"
}

func (c *SkipCommand) Description() string {
	return "Skip the current song"
}

func (c *SkipCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "skip",
		Description: "Skip the current song",
	}
}

func (c *SkipCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	voice.Skip(i.GuildID)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Skipped current song",
		},
	})
}
