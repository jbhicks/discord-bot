package imagegen

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type GenerationRequest struct {
	Prompt         string  `json:"prompt"`
	NegativePrompt string  `json:"negative_prompt,omitempty"`
	Steps          int     `json:"steps,omitempty"`
	Width          int     `json:"width,omitempty"`
	Height         int     `json:"height,omitempty"`
	CfgScale       float64 `json:"cfg_scale,omitempty"`
	SamplerName    string  `json:"sampler_name,omitempty"`
}

type GenerationResponse struct {
	Images     []string `json:"images"`
	Parameters struct {
		Prompt         string  `json:"prompt"`
		NegativePrompt string  `json:"negative_prompt"`
		Steps          int     `json:"steps"`
		Width          int     `json:"width"`
		Height         int     `json:"height"`
		CfgScale       float64 `json:"cfg_scale"`
		Seed           int64   `json:"seed"`
	} `json:"parameters"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) GenerateImage(req *GenerationRequest) (*GenerationResponse, error) {
	if req.Steps == 0 {
		req.Steps = 30
	}
	if req.Width == 0 {
		req.Width = 1024
	}
	if req.Height == 0 {
		req.Height = 1024
	}
	if req.CfgScale == 0 {
		req.CfgScale = 7.0
	}
	if req.SamplerName == "" {
		req.SamplerName = "DPM++ 2M"
	}
	if req.NegativePrompt == "" {
		req.NegativePrompt = "ugly, blurry, low quality, distorted, deformed, bad anatomy"
	}

	slog.Info("Generating image",
		"prompt", req.Prompt,
		"steps", req.Steps,
		"width", req.Width,
		"height", req.Height,
	)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/sdapi/v1/txt2img", c.baseURL)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var genResp GenerationResponse
	if err := json.Unmarshal(body, &genResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	slog.Info("Image generated successfully",
		"num_images", len(genResp.Images),
		"seed", genResp.Parameters.Seed,
	)

	return &genResp, nil
}

func (c *Client) DecodeImage(base64Str string) ([]byte, error) {
	imageData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 image: %w", err)
	}
	return imageData, nil
}

func (c *Client) HealthCheck() error {
	url := fmt.Sprintf("%s/sdapi/v1/sd-models", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	return nil
}
