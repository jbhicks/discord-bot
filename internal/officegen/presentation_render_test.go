package officegen

import (
	"os"
	"strings"
	"testing"
)

func TestPresentationSecondSlideRendering(t *testing.T) {
	gen := NewHTMLGenerator()
	pdfExport := NewPDFExporter()

	content := &PresentationContent{
		Title:    "Test Presentation",
		Subtitle: "Testing Second Slide Issue",
		Slides: []PresentationSlide{
			{
				Title: "First Content Slide",
				Bullets: []string{
					"First bullet point on first slide",
					"Second bullet point on first slide",
					"Third bullet point on first slide",
				},
			},
			{
				Title: "Second Content Slide",
				Bullets: []string{
					"First bullet point on second slide",
					"Second bullet point on second slide",
					"Third bullet point on second slide",
				},
			},
			{
				Title: "Third Content Slide",
				Bullets: []string{
					"First bullet point on third slide",
					"Second bullet point on third slide",
					"Third bullet point on third slide",
				},
			},
		},
	}

	largeFakeImage := make([]byte, 500000)
	for i := range largeFakeImage {
		largeFakeImage[i] = byte(i % 256)
	}
	images := [][]byte{largeFakeImage, largeFakeImage, largeFakeImage}

	html, err := gen.GeneratePresentationHTML(content, images)
	if err != nil {
		t.Fatalf("Failed to generate HTML: %v", err)
	}

	if err := os.WriteFile("/tmp/test_presentation.html", []byte(html), 0644); err != nil {
		t.Logf("Warning: Failed to save test HTML: %v", err)
	} else {
		t.Logf("Saved test HTML to /tmp/test_presentation.html")
	}

	t.Run("HTML contains all slide titles", func(t *testing.T) {
		if !strings.Contains(html, "First Content Slide") {
			t.Error("Missing first slide title")
		}
		if !strings.Contains(html, "Second Content Slide") {
			t.Error("Missing second slide title - this is the problematic slide")
		}
		if !strings.Contains(html, "Third Content Slide") {
			t.Error("Missing third slide title")
		}
	})

	t.Run("HTML contains all bullet points", func(t *testing.T) {
		if !strings.Contains(html, "First bullet point on second slide") {
			t.Error("Missing bullets from second slide")
		}
	})

	t.Run("HTML contains correct number of slides", func(t *testing.T) {
		contentSlideCount := strings.Count(html, `class="slide content-slide"`)
		if contentSlideCount != 3 {
			t.Errorf("Expected 3 content slides, got %d", contentSlideCount)
		}
	})

	t.Run("HTML contains images for all slides", func(t *testing.T) {
		imageCount := strings.Count(html, `class="slide-image"`)
		if imageCount != 3 {
			t.Errorf("Expected 3 slide images, got %d", imageCount)
		}
	})

	pdfData, err := pdfExport.ConvertHTMLToPDF(html)
	if err != nil {
		t.Fatalf("Failed to convert HTML to PDF: %v", err)
	}

	if err := os.WriteFile("/tmp/test_presentation.pdf", pdfData, 0644); err != nil {
		t.Fatalf("Failed to save test PDF: %v", err)
	}

	t.Logf("Generated PDF saved to /tmp/test_presentation.pdf (%d bytes)", len(pdfData))
	t.Logf("Please inspect the PDF to verify second slide renders correctly")

	if len(pdfData) == 0 {
		t.Fatal("PDF data is empty")
	}

	if len(pdfData) < 1000 {
		t.Errorf("PDF seems too small (%d bytes), might be corrupted", len(pdfData))
	}
}
