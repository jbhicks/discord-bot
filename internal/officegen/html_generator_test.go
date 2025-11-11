package officegen

import (
	"strings"
	"testing"
)

func TestGeneratePresentationHTML_WithTextAndImages(t *testing.T) {
	gen := NewHTMLGenerator()

	content := &PresentationContent{
		Title:    "Test Presentation",
		Subtitle: "A test subtitle",
		Slides: []PresentationSlide{
			{
				Title: "First Slide",
				Bullets: []string{
					"First bullet point",
					"Second bullet point",
					"Third bullet point",
				},
			},
			{
				Title: "Second Slide",
				Bullets: []string{
					"Another point",
					"One more point",
				},
			},
		},
	}

	fakeImage1 := []byte("fake-image-data-1")
	fakeImage2 := []byte("fake-image-data-2")
	images := [][]byte{fakeImage1, fakeImage2}

	html, err := gen.GeneratePresentationHTML(content, images)
	if err != nil {
		t.Fatalf("Failed to generate HTML: %v", err)
	}

	t.Run("Contains title slide", func(t *testing.T) {
		if !strings.Contains(html, "Test Presentation") {
			t.Error("HTML should contain presentation title")
		}
		if !strings.Contains(html, "A test subtitle") {
			t.Error("HTML should contain presentation subtitle")
		}
	})

	t.Run("Contains slide titles", func(t *testing.T) {
		if !strings.Contains(html, "First Slide") {
			t.Error("HTML should contain first slide title")
		}
		if !strings.Contains(html, "Second Slide") {
			t.Error("HTML should contain second slide title")
		}
	})

	t.Run("Contains slide bullet points", func(t *testing.T) {
		if !strings.Contains(html, "First bullet point") {
			t.Error("HTML should contain first bullet point")
		}
		if !strings.Contains(html, "Second bullet point") {
			t.Error("HTML should contain second bullet point")
		}
		if !strings.Contains(html, "Third bullet point") {
			t.Error("HTML should contain third bullet point")
		}
		if !strings.Contains(html, "Another point") {
			t.Error("HTML should contain bullet from second slide")
		}
		if !strings.Contains(html, "One more point") {
			t.Error("HTML should contain another bullet from second slide")
		}
	})

	t.Run("Contains images", func(t *testing.T) {
		if !strings.Contains(html, "data:image/png;base64,") {
			t.Error("HTML should contain base64 encoded images")
		}
		if strings.Count(html, "data:image/png;base64,") < 2 {
			t.Errorf("HTML should contain at least 2 images, found %d", strings.Count(html, "data:image/png;base64,"))
		}
	})

	t.Run("HTML structure is valid", func(t *testing.T) {
		if !strings.Contains(html, "<!DOCTYPE html>") {
			t.Error("HTML should start with DOCTYPE")
		}
		if !strings.Contains(html, "</html>") {
			t.Error("HTML should end with closing html tag")
		}
		if strings.Count(html, "<div class=\"slide content-slide\">") != 2 {
			t.Errorf("Expected 2 content slides, got %d", strings.Count(html, "<div class=\"slide content-slide\">"))
		}
	})
}

func TestGenerateDocumentHTML_WithTextAndImages(t *testing.T) {
	gen := NewHTMLGenerator()

	content := &DocumentContent{
		Title: "Test Document",
		Sections: []DocumentSection{
			{
				Heading: "Introduction",
				Content: "This is the introduction paragraph. It contains important information.",
			},
			{
				Heading: "Main Section",
				Content: "This is the main section with detailed content. Here is more information about the topic.",
			},
		},
	}

	fakeImage1 := []byte("fake-image-data-1")
	fakeImage2 := []byte("fake-image-data-2")
	images := [][]byte{fakeImage1, fakeImage2}

	html, err := gen.GenerateDocumentHTML(content, images)
	if err != nil {
		t.Fatalf("Failed to generate HTML: %v", err)
	}

	t.Run("Contains document title", func(t *testing.T) {
		if !strings.Contains(html, "Test Document") {
			t.Error("HTML should contain document title")
		}
	})

	t.Run("Contains section headings", func(t *testing.T) {
		if !strings.Contains(html, "Introduction") {
			t.Error("HTML should contain Introduction heading")
		}
		if !strings.Contains(html, "Main Section") {
			t.Error("HTML should contain Main Section heading")
		}
	})

	t.Run("Contains section content", func(t *testing.T) {
		if !strings.Contains(html, "This is the introduction paragraph") {
			t.Error("HTML should contain introduction content")
		}
		if !strings.Contains(html, "This is the main section with detailed content") {
			t.Error("HTML should contain main section content")
		}
	})

	t.Run("Contains images", func(t *testing.T) {
		if !strings.Contains(html, "data:image/png;base64,") {
			t.Error("HTML should contain base64 encoded images")
		}
	})

	t.Run("HTML structure is valid", func(t *testing.T) {
		if !strings.Contains(html, "<!DOCTYPE html>") {
			t.Error("HTML should start with DOCTYPE")
		}
		if !strings.Contains(html, "</html>") {
			t.Error("HTML should end with closing html tag")
		}
	})
}

func TestGeneratePresentationHTML_WithoutImages(t *testing.T) {
	gen := NewHTMLGenerator()

	content := &PresentationContent{
		Title:    "No Images Presentation",
		Subtitle: "Testing without images",
		Slides: []PresentationSlide{
			{
				Title: "Content Slide",
				Bullets: []string{
					"Point one",
					"Point two",
				},
			},
		},
	}

	html, err := gen.GeneratePresentationHTML(content, nil)
	if err != nil {
		t.Fatalf("Failed to generate HTML: %v", err)
	}

	if !strings.Contains(html, "Content Slide") {
		t.Error("HTML should contain slide title even without images")
	}
	if !strings.Contains(html, "Point one") {
		t.Error("HTML should contain bullet points even without images")
	}
}
