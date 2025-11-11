package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/josh/discord-bot/internal/voice"
)

type StopCommand struct{}

func (c *StopCommand) Name() string {
	return "stop"
}

func (c *StopCommand) Description() string {
	return "Stop playing and leave voice channel"
}

func (c *StopCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "stop",
		Description: "Stop playing and leave voice channel",
	}
}

func (c *StopCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	voice.StopPlaying(i.GuildID)
	voice.LeaveVoiceChannel(i.GuildID)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Record scratch!",
		},
	})
}
