package officegen

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type PDFExporter struct {
	htmlGen *HTMLGenerator
}

func NewPDFExporter() *PDFExporter {
	return &PDFExporter{
		htmlGen: NewHTMLGenerator(),
	}
}

func (pe *PDFExporter) ConvertHTMLToPDF(html string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var pdfBuf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				Do(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	slog.Info("PDF generated from HTML", "size", len(pdfBuf))

	if err := os.WriteFile("/tmp/generated_pdf_output.pdf", pdfBuf, 0644); err != nil {
		slog.Warn("Failed to save debug PDF", "error", err)
	} else {
		slog.Info("Saved debug PDF to /tmp/generated_pdf_output.pdf")
	}

	return pdfBuf, nil
}

func (pe *PDFExporter) CheckAvailable() bool {
	_, err := exec.LookPath("chromium")
	if err == nil {
		return true
	}
	_, err = exec.LookPath("chromium-browser")
	if err == nil {
		return true
	}
	_, err = exec.LookPath("google-chrome")
	if err == nil {
		return true
	}
	_, err = exec.LookPath("chrome")
	return err == nil
}
