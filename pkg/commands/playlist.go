package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/josh/discord-bot/internal/db"
	"github.com/josh/discord-bot/internal/voice"
)

type PlaylistCommand struct{}

func (c *PlaylistCommand) Name() string {
	return "playlist"
}

func (c *PlaylistCommand) Description() string {
	return "Manage playlists"
}

func (c *PlaylistCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "playlist",
		Description: "Manage playlists",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "create",
				Description: "Create a new playlist",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "name",
						Description: "Playlist name",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Add song to playlist",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "name",
						Description: "Playlist name",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "url",
						Description: "Song URL",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "play",
				Description: "Play a playlist",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "name",
						Description: "Playlist name",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "list",
				Description: "List your playlists",
			},
		},
	}
}

func (c *PlaylistCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	if len(data.Options) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid subcommand",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	sub := data.Options[0]
	userID := i.Member.User.ID
	switch sub.Name {
	case "create":
		name := sub.Options[0].StringValue()
		err := db.CreatePlaylist(userID, name)
		if err != nil {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error creating playlist: " + err.Error(),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Playlist '" + name + "' created",
			},
		})
	case "add":
		name := sub.Options[0].StringValue()
		url := sub.Options[1].StringValue()
		err := db.AddToPlaylist(userID, name, url)
		if err != nil {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error adding to playlist: " + err.Error(),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Added to playlist '" + name + "'",
			},
		})
	case "play":
		name := sub.Options[0].StringValue()
		songs, err := db.GetPlaylist(userID, name)
		if err != nil {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error loading playlist: " + err.Error(),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
		for _, url := range songs {
			if strings.TrimSpace(url) != "" {
				voice.AddToQueue(i.GuildID, url, url)
			}
		}
		if !voice.IsPlaying(i.GuildID) {
			voice.StartPlaying(i.GuildID)
		}
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Playing playlist '" + name + "'",
			},
		})
	case "list":
		names, err := db.ListPlaylists(userID)
		if err != nil {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error listing playlists: " + err.Error(),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
		content := "Your playlists:\n" + strings.Join(names, "\n")
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
	default:
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unknown subcommand",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}
