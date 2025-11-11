package commands

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/josh/discord-bot/internal/imagegen"
	"github.com/josh/discord-bot/internal/officegen"
)

type PDFCommand struct {
	docGen   *officegen.DocumentGenerator
	sheetGen *officegen.SpreadsheetGenerator
	presGen  *officegen.PresentationGenerator
}

func NewPDFCommand(llmURL string, imageGenURL string) *PDFCommand {
	llmClient := officegen.NewClient(llmURL)
	sdClient := imagegen.NewClient(imageGenURL)
	imageGen := officegen.NewImageGenerator(sdClient, llmClient)

	return &PDFCommand{
		docGen:   officegen.NewDocumentGenerator(llmClient, imageGen),
		sheetGen: officegen.NewSpreadsheetGenerator(llmClient, imageGen),
		presGen:  officegen.NewPresentationGenerator(llmClient, imageGen),
	}
}

func (c *PDFCommand) Name() string {
	return "pdf"
}

func (c *PDFCommand) Description() string {
	return "Generate PDF documents (reports, presentations, tables) with AI"
}

func (c *PDFCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.Name(),
		Description: c.Description(),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "type",
				Description: "Type of document to generate",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "Document/Report", Value: "document"},
					{Name: "Presentation/Slides", Value: "presentation"},
					{Name: "Spreadsheet/Table", Value: "spreadsheet"},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "Describe the content you want in the document",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "title",
				Description: "Custom title for the document (optional)",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "pages",
				Description: "Target number of pages/slides/sheets (0=auto, default: 0)",
				Required:    false,
			},
		},
	}
}

func (c *PDFCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options

	username := "Unknown"
	userID := "Unknown"
	if i.Member != nil && i.Member.User != nil {
		username = i.Member.User.Username
		userID = i.Member.User.ID
	}

	docType := options[0].StringValue()
	prompt := options[1].StringValue()

	var title string
	var pages int

	for _, opt := range options[2:] {
		switch opt.Name {
		case "title":
			title = opt.StringValue()
		case "pages":
			pages = int(opt.IntValue())
		}
	}

	slog.Info("PDF command received",
		"user", username,
		"user_id", userID,
		"guild_id", i.GuildID,
		"type", docType,
		"prompt", prompt,
		"title", title,
		"pages", pages,
	)

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "📄 Generating your PDF with AI images... This may take 1-3 minutes.",
		},
	}); err != nil {
		return err
	}

	var result *officegen.GeneratedDocument
	var err error

	switch docType {
	case "document":
		result, err = c.docGen.Generate(&officegen.DocumentRequest{
			Prompt:      prompt,
			Title:       title,
			TargetPages: pages,
		})
	case "spreadsheet":
		result, err = c.sheetGen.Generate(&officegen.SpreadsheetRequest{
			Prompt:      prompt,
			Title:       title,
			TargetPages: pages,
		})
	case "presentation":
		result, err = c.presGen.Generate(&officegen.PresentationRequest{
			Prompt:       prompt,
			Title:        title,
			TargetSlides: pages,
		})
	default:
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: strPtr(fmt.Sprintf("❌ Unknown document type: %s", docType)),
		})
		return err
	}

	if err != nil {
		slog.Error("Failed to generate PDF", "error", err, "type", docType)
		_, editErr := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: strPtr(fmt.Sprintf("❌ Failed to generate PDF: %v", err)),
		})
		if editErr != nil {
			return editErr
		}
		return err
	}

	files := []*discordgo.File{
		{
			Name:        result.Filename,
			ContentType: "application/pdf",
			Reader:      bytes.NewReader(result.Data),
		},
	}

	content := fmt.Sprintf("✨ **PDF Generated**\n**Type:** %s\n**Filename:** %s\n**Size:** %.2f KB",
		docType, result.Filename, float64(len(result.Data))/1024)

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
		Files:   files,
	})

	if err != nil {
		slog.Error("Failed to send PDF", "error", err)
		return err
	}

	slog.Info("PDF sent successfully",
		"type", docType,
		"filename", result.Filename,
		"size", len(result.Data),
	)

	return nil
}
