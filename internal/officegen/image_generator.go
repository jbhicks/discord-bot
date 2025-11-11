package officegen

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/josh/discord-bot/internal/imagegen"
)

type ImageGenerator struct {
	sdClient  *imagegen.Client
	llmClient *Client
}

func NewImageGenerator(sdClient *imagegen.Client, llmClient *Client) *ImageGenerator {
	return &ImageGenerator{
		sdClient:  sdClient,
		llmClient: llmClient,
	}
}

func (ig *ImageGenerator) GenerateImagePrompts(content any, count int) ([]string, error) {
	if count == 0 {
		count = 2
	}

	contentJSON, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`Based on this content, generate %d specific, detailed image generation prompts for Stable Diffusion.

Requirements:
- Each prompt should be 20-40 words
- Focus on visual, concrete elements (no abstract concepts)
- Include style descriptors (e.g., "digital art", "photorealistic", "illustration")
- Make each prompt unique and relevant to different parts of the content

IMPORTANT: Return ONLY a valid JSON array of strings with no markdown formatting or code blocks.

Format: ["prompt1", "prompt2", ...]

Content: %s`, count, string(contentJSON))

	response, err := ig.llmClient.GenerateText(prompt)
	if err != nil {
		return nil, err
	}

	cleanedResponse := cleanJSONResponse(response)

	var prompts []string
	if err := json.Unmarshal([]byte(cleanedResponse), &prompts); err != nil {
		slog.Warn("Failed to parse image prompts, creating default", "error", err, "response", cleanedResponse[:min(200, len(cleanedResponse))])
		prompts = []string{
			"professional illustration related to the topic",
			"diagram or visualization of key concepts",
		}
	}

	if len(prompts) > count {
		prompts = prompts[:count]
	}

	return prompts, nil
}

func (ig *ImageGenerator) GenerateImages(prompts []string) ([][]byte, error) {
	if err := ig.sdClient.HealthCheck(); err != nil {
		slog.Warn("SD WebUI not available, skipping image generation", "error", err)
		return nil, fmt.Errorf("image generation service unavailable: %w", err)
	}

	var images [][]byte

	for i, prompt := range prompts {
		slog.Info("Generating image", "index", i+1, "total", len(prompts), "prompt", prompt)

		req := &imagegen.GenerationRequest{
			Prompt: prompt,
			Width:  768,
			Height: 512,
			Steps:  20,
		}

		resp, err := ig.sdClient.GenerateImage(req)
		if err != nil {
			slog.Warn("Failed to generate image", "prompt", prompt, "error", err)
			continue
		}

		if len(resp.Images) == 0 {
			slog.Warn("No image data returned", "prompt", prompt)
			continue
		}

		imageData, err := ig.sdClient.DecodeImage(resp.Images[0])
		if err != nil {
			slog.Warn("Failed to decode image", "prompt", prompt, "error", err)
			continue
		}

		images = append(images, imageData)
		slog.Info("Image generated successfully", "index", i+1, "size", len(imageData))
	}

	return images, nil
}
