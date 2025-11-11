package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/josh/discord-bot/internal/voice"
)

type PlayCommand struct{}

func (c *PlayCommand) Name() string {
	return "play"
}

func (c *PlayCommand) Description() string {
	return "Play a song from URL"
}

func (c *PlayCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "play",
		Description: "Play a song from URL",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "url",
				Description: "URL of the song to play",
				Required:    true,
			},
		},
	}
}

func (c *PlayCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	url := i.ApplicationCommandData().Options[0].StringValue()

	// Check if user is in voice channel
	vs, voiceErr := s.State.VoiceState(i.GuildID, i.Member.User.ID)
	if voiceErr != nil || vs == nil || vs.ChannelID == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You must be in a voice channel to use this command!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	// Join voice channel
	err := voice.JoinVoiceChannel(s, i.GuildID, vs.ChannelID)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to join voice channel!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	// Respond first
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Now playing: " + url,
		},
	})
	if err != nil {
		return err
	}

	// Add to queue
	if strings.HasSuffix(url, ".mp3") {
		voice.AddToQueue(i.GuildID, url, url)
		if !voice.IsPlaying(i.GuildID) {
			voice.StartPlaying(i.GuildID)
		}
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Added to queue: " + url,
			},
		})
	} else {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Only MP3 URLs are supported",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}
