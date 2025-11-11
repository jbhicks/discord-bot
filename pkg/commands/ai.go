package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type AICommand struct{}

func (c *AICommand) Name() string {
	return "ai"
}

func (c *AICommand) Description() string {
	return "Ask the AI a question"
}

func (c *AICommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.Name(),
		Description: c.Description(),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "The prompt to send to the AI",
				Required:    true,
			},
		},
	}
}

func (c *AICommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	prompt := i.ApplicationCommandData().Options[0].StringValue()

	username := "Unknown"
	userID := "Unknown"
	if i.Member != nil && i.Member.User != nil {
		username = i.Member.User.Username
		userID = i.Member.User.ID
	}

	slog.Info("AI command received",
		"user", username,
		"user_id", userID,
		"guild_id", i.GuildID,
		"prompt", prompt,
	)

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		return err
	}

	requestBody := map[string]interface{}{
		"model": "llama",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8081/v1/chat/completions", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return fmt.Errorf("invalid response from AI server")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response structure")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid message structure")
	}

	content, ok := message["content"].(string)
	if !ok {
		return fmt.Errorf("invalid content")
	}

	const maxLen = 2000

	chunks := make([]string, 0)
	for len(content) > 0 {
		end := maxLen
		if end > len(content) {
			end = len(content)
		}
		chunks = append(chunks, content[:end])
		content = content[end:]
	}

	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &chunks[0],
	}); err != nil {
		return err
	}

	for _, chunk := range chunks[1:] {
		if _, err := s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: chunk,
		}); err != nil {
			return err
		}
	}

	return nil
}
