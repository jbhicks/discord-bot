package commands

import (
	"net/url"

	"github.com/bwmarrin/discordgo"
)

type SearchCommand struct{}

func (c *SearchCommand) Name() string {
	return "search"
}

func (c *SearchCommand) Description() string {
	return "Search for songs on YouTube"
}

func (c *SearchCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "search",
		Description: "Search for songs on YouTube",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "Search query",
				Required:    true,
			},
		},
	}
}

func (c *SearchCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	query := i.ApplicationCommandData().Options[0].StringValue()
	searchURL := "https://www.youtube.com/results?search_query=" + url.QueryEscape(query)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Search results: " + searchURL,
		},
	})
}
