package officegen

import (
	"fmt"
	"log/slog"
)

type PresentationGenerator struct {
	llmClient *Client
	imageGen  *ImageGenerator
	htmlGen   *HTMLGenerator
	pdfExport *PDFExporter
}

func NewPresentationGenerator(llmClient *Client, imageGen *ImageGenerator) *PresentationGenerator {
	return &PresentationGenerator{
		llmClient: llmClient,
		imageGen:  imageGen,
		htmlGen:   NewHTMLGenerator(),
		pdfExport: NewPDFExporter(),
	}
}

func (pg *PresentationGenerator) Generate(req *PresentationRequest) (*GeneratedDocument, error) {
	slog.Info("Generating presentation", "prompt", req.Prompt, "target_slides", req.TargetSlides)

	content, err := pg.llmClient.GeneratePresentation(req.Prompt, req.TargetSlides)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if req.Title != "" {
		content.Title = req.Title
	}

	slog.Info("Before filtering slides", "count", len(content.Slides))
	for i, slide := range content.Slides {
		slog.Info("Slide before filtering", "index", i, "title", slide.Title, "bullets", len(slide.Bullets))
	}

	var validSlides []PresentationSlide
	for i, slide := range content.Slides {
		if slide.Title == "" && len(slide.Bullets) == 0 {
			slog.Warn("Skipping empty slide", "index", i)
			continue
		}
		validSlides = append(validSlides, slide)
	}
	content.Slides = validSlides

	slog.Info("After filtering slides", "count", len(content.Slides))
	for i, slide := range content.Slides {
		slog.Info("Slide after filtering", "index", i, "title", slide.Title, "bullets", len(slide.Bullets))
	}

	var images [][]byte
	if pg.imageGen != nil {
		imageCount := min(len(content.Slides), 5)

		slog.Info("Generating image prompts for presentation", "count", imageCount)
		prompts, err := pg.imageGen.GenerateImagePrompts(content, imageCount)
		if err != nil {
			slog.Warn("Failed to generate image prompts", "error", err)
		} else {
			slog.Info("Generating images for presentation", "prompts", prompts)
			images, err = pg.imageGen.GenerateImages(prompts)
			if err != nil {
				slog.Warn("Failed to generate images", "error", err)
			} else {
				slog.Info("Generated images for presentation", "count", len(images))
			}
		}
	}

	html, err := pg.htmlGen.GeneratePresentationHTML(content, images)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML: %w", err)
	}

	pdfData, err := pg.pdfExport.ConvertHTMLToPDF(html)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to PDF: %w", err)
	}

	filename := sanitizeFilename(content.Title) + ".pdf"
	slog.Info("Presentation generated", "filename", filename, "size", len(pdfData))

	return &GeneratedDocument{
		Data:     pdfData,
		Filename: filename,
	}, nil
}
