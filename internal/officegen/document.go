package officegen

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
)

type DocumentGenerator struct {
	llmClient *Client
	imageGen  *ImageGenerator
	htmlGen   *HTMLGenerator
	pdfExport *PDFExporter
}

func NewDocumentGenerator(llmClient *Client, imageGen *ImageGenerator) *DocumentGenerator {
	return &DocumentGenerator{
		llmClient: llmClient,
		imageGen:  imageGen,
		htmlGen:   NewHTMLGenerator(),
		pdfExport: NewPDFExporter(),
	}
}

func (dg *DocumentGenerator) Generate(req *DocumentRequest) (*GeneratedDocument, error) {
	slog.Info("Generating document", "prompt", req.Prompt, "target_pages", req.TargetPages)

	content, err := dg.llmClient.GenerateDocument(req.Prompt, req.TargetPages, true)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if req.Title != "" {
		content.Title = req.Title
	}

	var images [][]byte
	if dg.imageGen != nil {
		imageCount := min(3, len(content.Sections)/2+1)

		slog.Info("Generating image prompts", "count", imageCount)
		prompts, err := dg.imageGen.GenerateImagePrompts(content, imageCount)
		if err != nil {
			slog.Warn("Failed to generate image prompts", "error", err)
		} else {
			slog.Info("Generating images", "prompts", prompts)
			images, err = dg.imageGen.GenerateImages(prompts)
			if err != nil {
				slog.Warn("Failed to generate images", "error", err)
			} else {
				slog.Info("Generated images", "count", len(images))
			}
		}
	}

	html, err := dg.htmlGen.GenerateDocumentHTML(content, images)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML: %w", err)
	}

	pdfData, err := dg.pdfExport.ConvertHTMLToPDF(html)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to PDF: %w", err)
	}

	filename := sanitizeFilename(content.Title) + ".pdf"
	slog.Info("Document generated", "filename", filename, "size", len(pdfData))

	return &GeneratedDocument{
		Data:     pdfData,
		Filename: filename,
	}, nil
}

func sanitizeFilename(name string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9\-_\s]`)
	clean := reg.ReplaceAllString(name, "")
	clean = strings.TrimSpace(clean)
	clean = regexp.MustCompile(`\s+`).ReplaceAllString(clean, "_")

	if clean == "" {
		clean = "document"
	}

	if len(clean) > 100 {
		clean = clean[:100]
	}

	return clean
}
