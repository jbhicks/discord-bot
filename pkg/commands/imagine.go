package commands

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/josh/discord-bot/internal/imagegen"
)

type ImagineCommand struct {
	client *imagegen.Client
}

func NewImagineCommand(imageGenURL string) *ImagineCommand {
	return &ImagineCommand{
		client: imagegen.NewClient(imageGenURL),
	}
}

func (c *ImagineCommand) Name() string {
	return "imagine"
}

func (c *ImagineCommand) Description() string {
	return "Generate an image from a text prompt using AI"
}

func (c *ImagineCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.Name(),
		Description: c.Description(),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "Describe the image you want to generate",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "width",
				Description: "Image width (default: 1024, max: 2048)",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "height",
				Description: "Image height (default: 1024, max: 2048)",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "steps",
				Description: "Number of generation steps (default: 30, max: 50)",
				Required:    false,
			},
		},
	}
}

func (c *ImagineCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
	prompt := options[0].StringValue()

	username := "Unknown"
	userID := "Unknown"
	if i.Member != nil && i.Member.User != nil {
		username = i.Member.User.Username
		userID = i.Member.User.ID
	}

	slog.Info("Imagine command received",
		"user", username,
		"user_id", userID,
		"guild_id", i.GuildID,
		"prompt", prompt,
	)

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🎨 Generating your image... This may take 30-60 seconds.",
		},
	}); err != nil {
		return err
	}

	req := &imagegen.GenerationRequest{
		Prompt: prompt,
	}

	for _, opt := range options[1:] {
		switch opt.Name {
		case "width":
			req.Width = int(opt.IntValue())
		case "height":
			req.Height = int(opt.IntValue())
		case "steps":
			req.Steps = int(opt.IntValue())
		}
	}

	resp, err := c.client.GenerateImage(req)
	if err != nil {
		slog.Error("Failed to generate image", "error", err)
		_, editErr := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: strPtr(fmt.Sprintf("❌ Failed to generate image: %v", err)),
		})
		if editErr != nil {
			return editErr
		}
		return err
	}

	if len(resp.Images) == 0 {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: strPtr("❌ No image was generated"),
		})
		return err
	}

	imageData, err := c.client.DecodeImage(resp.Images[0])
	if err != nil {
		slog.Error("Failed to decode image", "error", err)
		_, editErr := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: strPtr(fmt.Sprintf("❌ Failed to decode image: %v", err)),
		})
		if editErr != nil {
			return editErr
		}
		return err
	}

	file := &discordgo.File{
		Name:        "generated_image.png",
		ContentType: "image/png",
		Reader:      bytes.NewReader(imageData),
	}

	content := fmt.Sprintf("✨ **Generated Image**\n**Prompt:** %s\n**Size:** %dx%d | **Steps:** %d | **Seed:** %d",
		resp.Parameters.Prompt,
		resp.Parameters.Width,
		resp.Parameters.Height,
		resp.Parameters.Steps,
		resp.Parameters.Seed,
	)

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
		Files:   []*discordgo.File{file},
	})

	if err != nil {
		slog.Error("Failed to send image", "error", err)
		return err
	}

	slog.Info("Image sent successfully",
		"user", username,
		"seed", resp.Parameters.Seed,
	)

	return nil
}

func strPtr(s string) *string {
	return &s
}
