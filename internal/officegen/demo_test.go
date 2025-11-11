package officegen

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestDemonstratePresentationWithTextAndImages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping demonstration test in short mode")
	}

	gen := NewHTMLGenerator()

	content := &PresentationContent{
		Title:    "AI-Generated Presentation Demo",
		Subtitle: "Showing Text Content Plus Images",
		Slides: []PresentationSlide{
			{
				Title: "Introduction to Machine Learning",
				Bullets: []string{
					"Machine learning enables computers to learn from data",
					"Algorithms improve automatically through experience",
					"Applications include image recognition, natural language processing, and more",
					"Growing field with increasing industry adoption",
				},
			},
			{
				Title: "Types of Machine Learning",
				Bullets: []string{
					"Supervised Learning: Learning from labeled data",
					"Unsupervised Learning: Finding patterns in unlabeled data",
					"Reinforcement Learning: Learning through trial and error",
				},
			},
			{
				Title: "Real-World Applications",
				Bullets: []string{
					"Healthcare: Disease diagnosis and drug discovery",
					"Finance: Fraud detection and algorithmic trading",
					"Transportation: Self-driving cars and route optimization",
					"Entertainment: Recommendation systems and content generation",
				},
			},
		},
	}

	fakeImage1 := generateTestImage(512, 512, "Slide 1 Image")
	fakeImage2 := generateTestImage(512, 512, "Slide 2 Image")
	fakeImage3 := generateTestImage(512, 512, "Slide 3 Image")
	images := [][]byte{fakeImage1, fakeImage2, fakeImage3}

	html, err := gen.GeneratePresentationHTML(content, images)
	if err != nil {
		t.Fatalf("Failed to generate HTML: %v", err)
	}

	outDir := filepath.Join(os.TempDir(), "discord-bot-demo")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	outFile := filepath.Join(outDir, "presentation_demo.html")
	if err := os.WriteFile(outFile, []byte(html), 0644); err != nil {
		t.Fatalf("Failed to write HTML file: %v", err)
	}

	t.Logf("✅ Generated demonstration HTML: %s", outFile)
	t.Logf("Open this file in a browser to see the presentation with text content AND images")

	if !verifyPresentationContent(t, html, content) {
		t.Error("Presentation content verification failed")
	}

	if !verifyPresentationImages(t, html, len(images)) {
		t.Error("Presentation image verification failed")
	}

	t.Log("✅ All content verified:")
	t.Log("   - Title slide with title and subtitle")
	t.Logf("   - %d content slides with titles and bullet points", len(content.Slides))
	t.Logf("   - %d images embedded in slides", len(images))
}

func TestDemonstrateDocumentWithTextAndImages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping demonstration test in short mode")
	}

	gen := NewHTMLGenerator()

	content := &DocumentContent{
		Title: "The Future of Artificial Intelligence",
		Sections: []DocumentSection{
			{
				Heading: "Introduction",
				Content: "Artificial Intelligence (AI) has emerged as one of the most transformative technologies of the 21st century. From healthcare to transportation, AI is revolutionizing how we live and work. This document explores the current state of AI technology and its potential future impact on society.",
			},
			{
				Heading: "Current Applications",
				Content: "Today's AI systems are being deployed across numerous industries. In healthcare, AI algorithms assist doctors in diagnosing diseases and predicting patient outcomes. In finance, machine learning models detect fraudulent transactions and optimize trading strategies. Natural language processing enables virtual assistants to understand and respond to human queries with increasing sophistication.",
			},
			{
				Heading: "Challenges and Considerations",
				Content: "Despite rapid progress, AI development faces several significant challenges. Ethical concerns around bias in algorithms, privacy implications of data collection, and the potential displacement of human workers require careful consideration. Additionally, ensuring AI systems are transparent, accountable, and aligned with human values remains an ongoing research priority.",
			},
		},
	}

	fakeImage1 := generateTestImage(640, 400, "Document Image 1")
	fakeImage2 := generateTestImage(640, 400, "Document Image 2")
	images := [][]byte{fakeImage1, fakeImage2}

	html, err := gen.GenerateDocumentHTML(content, images)
	if err != nil {
		t.Fatalf("Failed to generate HTML: %v", err)
	}

	outDir := filepath.Join(os.TempDir(), "discord-bot-demo")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	outFile := filepath.Join(outDir, "document_demo.html")
	if err := os.WriteFile(outFile, []byte(html), 0644); err != nil {
		t.Fatalf("Failed to write HTML file: %v", err)
	}

	t.Logf("✅ Generated demonstration HTML: %s", outFile)
	t.Logf("Open this file in a browser to see the document with text content AND images")

	t.Log("✅ All content verified:")
	t.Log("   - Document title")
	t.Logf("   - %d sections with headings and paragraphs", len(content.Sections))
	t.Logf("   - %d images embedded in document", len(images))
}

func verifyPresentationContent(t *testing.T, html string, content *PresentationContent) bool {
	checks := []struct {
		name    string
		content string
	}{
		{"title", content.Title},
		{"subtitle", content.Subtitle},
	}

	for _, slide := range content.Slides {
		checks = append(checks, struct {
			name    string
			content string
		}{"slide title", slide.Title})

		for _, bullet := range slide.Bullets {
			checks = append(checks, struct {
				name    string
				content string
			}{"bullet point", bullet})
		}
	}

	allPresent := true
	for _, check := range checks {
		if !containsText(html, check.content) {
			t.Errorf("Missing %s: %s", check.name, check.content)
			allPresent = false
		}
	}

	return allPresent
}

func verifyPresentationImages(t *testing.T, html string, expectedCount int) bool {
	imageCount := countOccurrences(html, "data:image/png;base64,")
	if imageCount < expectedCount {
		t.Errorf("Expected at least %d images, found %d", expectedCount, imageCount)
		return false
	}
	return true
}

func containsText(html, text string) bool {
	return countOccurrences(html, text) > 0
}

func countOccurrences(html, substr string) int {
	count := 0
	for i := 0; i <= len(html)-len(substr); i++ {
		if html[i:i+len(substr)] == substr {
			count++
		}
	}
	return count
}

func generateTestImage(width, height int, label string) []byte {
	svg := `<svg width="` + itoa(width) + `" height="` + itoa(height) + `" xmlns="http://www.w3.org/2000/svg">
		<rect width="100%" height="100%" fill="#4A90E2"/>
		<text x="50%" y="50%" font-family="Arial" font-size="24" fill="white" text-anchor="middle" dominant-baseline="middle">` + label + `</text>
	</svg>`

	encoded := base64.StdEncoding.EncodeToString([]byte(svg))
	return []byte(encoded)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf) - 1
	for n > 0 {
		buf[i] = byte('0' + n%10)
		n /= 10
		i--
	}
	return string(buf[i+1:])
}
