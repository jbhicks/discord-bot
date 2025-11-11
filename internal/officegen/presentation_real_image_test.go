package officegen

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type SDResponse struct {
	Images []string `json:"images"`
}

type SDRequest struct {
	Prompt string `json:"prompt"`
	Steps  int    `json:"steps"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func generateRealImageFromSD(prompt string) ([]byte, error) {
	reqData := SDRequest{
		Prompt: prompt,
		Steps:  20,
		Width:  768,
		Height: 512,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Post(
		"http://localhost:7860/sdapi/v1/txt2img",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sdResp SDResponse
	if err := json.Unmarshal(body, &sdResp); err != nil {
		return nil, err
	}

	if len(sdResp.Images) == 0 {
		return nil, fmt.Errorf("no images returned from SD")
	}

	imageData, err := base64.StdEncoding.DecodeString(sdResp.Images[0])
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func TestPresentationWithRealGeneratedImages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real image generation test in short mode")
	}

	t.Log("Generating 3 test images with Stable Diffusion (this may take a minute)...")

	images := make([][]byte, 3)
	prompts := []string{
		"a serene mountain landscape at sunset",
		"abstract geometric patterns in blue and purple",
		"a futuristic city skyline at night",
	}

	for i, prompt := range prompts {
		t.Logf("Generating image %d/3: %s", i+1, prompt)
		img, err := generateRealImageFromSD(prompt)
		if err != nil {
			t.Fatalf("Failed to generate image %d: %v", i+1, err)
		}
		images[i] = img
		t.Logf("Generated image %d: %d bytes", i+1, len(img))
	}

	gen := NewHTMLGenerator()
	pdfExport := NewPDFExporter()

	content := &PresentationContent{
		Title:    "Real Image Test Presentation",
		Subtitle: "Testing with Stable Diffusion Generated Images (768x512)",
		Slides: []PresentationSlide{
			{
				Title: "First Content Slide (Landscape)",
				Bullets: []string{
					"This slide contains a landscape image",
					"Image is 768x512 pixels from Stable Diffusion",
					"Testing if second slide is blank",
				},
			},
			{
				Title: "Second Content Slide (Abstract)",
				Bullets: []string{
					"THIS IS THE CRITICAL SECOND SLIDE",
					"If this slide is blank, the bug is confirmed",
					"Abstract geometric patterns image below",
				},
			},
			{
				Title: "Third Content Slide (Cityscape)",
				Bullets: []string{
					"This is the third content slide",
					"Contains a futuristic city image",
					"Testing overall layout consistency",
				},
			},
		},
	}

	html, err := gen.GeneratePresentationHTML(content, images)
	if err != nil {
		t.Fatalf("Failed to generate HTML: %v", err)
	}

	if err := os.WriteFile("/tmp/test_real_images.html", []byte(html), 0644); err != nil {
		t.Logf("Warning: Failed to save test HTML: %v", err)
	} else {
		t.Logf("Saved test HTML to /tmp/test_real_images.html")
	}

	pdfData, err := pdfExport.ConvertHTMLToPDF(html)
	if err != nil {
		t.Fatalf("Failed to convert HTML to PDF: %v", err)
	}

	outputPath := "/tmp/test_real_images.pdf"
	if err := os.WriteFile(outputPath, pdfData, 0644); err != nil {
		t.Fatalf("Failed to save test PDF: %v", err)
	}

	t.Logf("✅ Generated PDF with real SD images saved to %s (%d bytes)", outputPath, len(pdfData))
	t.Logf("📄 Open this file to verify:")
	t.Logf("   - Second slide is NOT blank")
	t.Logf("   - All slides have properly positioned text and images")
	t.Logf("   - Images are not overlapping text")

	if len(pdfData) < 10000 {
		t.Errorf("PDF seems too small (%d bytes), might be corrupted", len(pdfData))
	}
}
