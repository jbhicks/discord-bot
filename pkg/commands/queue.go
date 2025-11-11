package commands

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/josh/discord-bot/internal/voice"
)

type QueueCommand struct{}

func (c *QueueCommand) Name() string {
	return "queue"
}

func (c *QueueCommand) Description() string {
	return "Manage the song queue"
}

func (c *QueueCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "queue",
		Description: "Manage the song queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Add a song to the queue",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "url",
						Description: "URL of the song",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Remove a song from the queue",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "index",
						Description: "Index of the song to remove (0-based)",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "view",
				Description: "View the current queue",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "shuffle",
				Description: "Shuffle the queue",
			},
		},
	}
}

func (c *QueueCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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
	switch sub.Name {
	case "add":
		url := sub.Options[0].StringValue()
		voice.AddToQueue(i.GuildID, url, url) // title as url for now
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Added to queue: " + url,
			},
		})
	case "remove":
		index := int(sub.Options[0].IntValue())
		err := voice.RemoveFromQueue(i.GuildID, index)
		if err != nil {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error: " + err.Error(),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Removed song at index " + strconv.Itoa(index),
			},
		})
	case "view":
		q := voice.ViewQueue(i.GuildID)
		if len(q) == 0 {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Queue is empty",
				},
			})
		}
		content := "Current queue:\n"
		for idx, song := range q {
			content += fmt.Sprintf("%d. %s\n", idx, song.Title)
		}
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
	case "shuffle":
		voice.ShuffleQueue(i.GuildID)
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Queue shuffled",
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
