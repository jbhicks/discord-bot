package stocknews

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type MarketAuxClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type marketAuxResponse struct {
	Data []marketAuxArticle `json:"data"`
	Meta struct {
		Found int `json:"found"`
		Page  int `json:"page"`
	} `json:"meta"`
}

type marketAuxArticle struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	URL         string             `json:"url"`
	PublishedAt time.Time          `json:"published_at"`
	Source      string             `json:"source"`
	Entities    []marketAuxEntity  `json:"entities"`
	Sentiment   marketAuxSentiment `json:"sentiment"`
}

type marketAuxEntity struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type marketAuxSentiment struct {
	Polarity string  `json:"polarity"`
	Score    float64 `json:"score"`
}

func NewMarketAuxClient(apiKey string) *MarketAuxClient {
	return &MarketAuxClient{
		apiKey:  apiKey,
		baseURL: "https://api.marketaux.com/v1",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *MarketAuxClient) GetNews(ticker string, days int) ([]NewsItem, error) {
	params := url.Values{}
	params.Add("api_token", c.apiKey)
	params.Add("symbols", ticker)
	params.Add("language", "en")
	params.Add("limit", "10")

	if days > 0 {
		fromDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
		params.Add("published_after", fromDate)
	}

	endpoint := fmt.Sprintf("%s/news/all?%s", c.baseURL, params.Encode())

	slog.Info("Fetching MarketAux news", "ticker", ticker, "days", days)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %w", err)
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

	var apiResp marketAuxResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	news := make([]NewsItem, 0, len(apiResp.Data))
	for _, article := range apiResp.Data {
		sentiment := article.Sentiment.Polarity
		if sentiment == "" {
			sentiment = "neutral"
		}

		news = append(news, NewsItem{
			Title:       article.Title,
			Description: article.Description,
			URL:         article.URL,
			Date:        article.PublishedAt,
			Source:      article.Source,
			Sentiment:   sentiment,
		})
	}

	slog.Info("MarketAux news fetched", "ticker", ticker, "count", len(news))

	return news, nil
}

func (c *MarketAuxClient) GetTrendingNews(days int) ([]NewsItem, error) {
	params := url.Values{}
	params.Add("api_token", c.apiKey)
	params.Add("language", "en")
	params.Add("limit", "100")

	if days > 0 {
		fromDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
		params.Add("published_after", fromDate)
	}

	endpoint := fmt.Sprintf("%s/news/all?%s", c.baseURL, params.Encode())

	slog.Info("Fetching MarketAux trending news", "days", days)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %w", err)
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

	var apiResp marketAuxResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	news := make([]NewsItem, 0, len(apiResp.Data))
	for _, article := range apiResp.Data {
		sentiment := article.Sentiment.Polarity
		if sentiment == "" {
			sentiment = "neutral"
		}

		news = append(news, NewsItem{
			Title:       article.Title,
			Description: article.Description,
			URL:         article.URL,
			Date:        article.PublishedAt,
			Source:      article.Source,
			Sentiment:   sentiment,
		})
	}

	slog.Info("MarketAux trending news fetched", "count", len(news))

	return news, nil
}

func (c *MarketAuxClient) HealthCheck() error {
	params := url.Values{}
	params.Add("api_token", c.apiKey)
	params.Add("limit", "1")

	endpoint := fmt.Sprintf("%s/news/all?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	return nil
}
