package officegen

import (
	"fmt"
	"log/slog"
)

type SpreadsheetGenerator struct {
	llmClient *Client
	imageGen  *ImageGenerator
	htmlGen   *HTMLGenerator
	pdfExport *PDFExporter
}

func NewSpreadsheetGenerator(llmClient *Client, imageGen *ImageGenerator) *SpreadsheetGenerator {
	return &SpreadsheetGenerator{
		llmClient: llmClient,
		imageGen:  imageGen,
		htmlGen:   NewHTMLGenerator(),
		pdfExport: NewPDFExporter(),
	}
}

func (sg *SpreadsheetGenerator) Generate(req *SpreadsheetRequest) (*GeneratedDocument, error) {
	slog.Info("Generating spreadsheet", "prompt", req.Prompt, "target_pages", req.TargetPages)

	content, err := sg.llmClient.GenerateSpreadsheet(req.Prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if req.Title != "" {
		content.Title = req.Title
	}

	html, err := sg.htmlGen.GenerateSpreadsheetHTML(content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML: %w", err)
	}

	pdfData, err := sg.pdfExport.ConvertHTMLToPDF(html)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to PDF: %w", err)
	}

	filename := sanitizeFilename(content.Title) + ".pdf"
	slog.Info("Spreadsheet generated", "filename", filename, "size", len(pdfData))

	return &GeneratedDocument{
		Data:     pdfData,
		Filename: filename,
	}, nil
}
