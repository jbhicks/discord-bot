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

type AlphaVantageClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type alphaVantageResponse struct {
	Feed                     []alphaVantageFeed `json:"feed"`
	SentimentScoreDefinition string             `json:"sentiment_score_definition"`
}

type alphaVantageFeed struct {
	Title                 string                        `json:"title"`
	URL                   string                        `json:"url"`
	TimePublished         string                        `json:"time_published"`
	Authors               []string                      `json:"authors"`
	Summary               string                        `json:"summary"`
	Source                string                        `json:"source"`
	OverallSentimentScore float64                       `json:"overall_sentiment_score"`
	OverallSentimentLabel string                        `json:"overall_sentiment_label"`
	TickerSentiment       []alphaVantageTickerSentiment `json:"ticker_sentiment"`
}

type alphaVantageTickerSentiment struct {
	Ticker               string `json:"ticker"`
	RelevanceScore       string `json:"relevance_score"`
	TickerSentimentScore string `json:"ticker_sentiment_score"`
	TickerSentimentLabel string `json:"ticker_sentiment_label"`
}

func NewAlphaVantageClient(apiKey string) *AlphaVantageClient {
	return &AlphaVantageClient{
		apiKey:  apiKey,
		baseURL: "https://www.alphavantage.co/query",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *AlphaVantageClient) GetNews(ticker string, days int) ([]NewsItem, error) {
	params := url.Values{}
	params.Add("function", "NEWS_SENTIMENT")
	params.Add("tickers", ticker)
	params.Add("apikey", c.apiKey)
	params.Add("limit", "50")

	if days > 0 {
		fromDate := time.Now().AddDate(0, 0, -days).Format("20060102T0000")
		params.Add("time_from", fromDate)
	}

	endpoint := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	slog.Info("Fetching Alpha Vantage news", "ticker", ticker, "days", days)

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

	var apiResp alphaVantageResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	news := make([]NewsItem, 0, len(apiResp.Feed))
	for _, feed := range apiResp.Feed {
		publishedAt, err := time.Parse("20060102T150405", feed.TimePublished)
		if err != nil {
			publishedAt = time.Now()
		}

		sentiment := feed.OverallSentimentLabel
		if sentiment == "" {
			sentiment = "neutral"
		}

		news = append(news, NewsItem{
			Title:       feed.Title,
			Description: feed.Summary,
			URL:         feed.URL,
			Date:        publishedAt,
			Source:      feed.Source,
			Sentiment:   sentiment,
		})
	}

	slog.Info("Alpha Vantage news fetched", "ticker", ticker, "count", len(news))

	return news, nil
}

func (c *AlphaVantageClient) GetTrendingNews(days int) ([]NewsItem, error) {
	params := url.Values{}
	params.Add("function", "NEWS_SENTIMENT")
	params.Add("apikey", c.apiKey)
	params.Add("limit", "100")

	if days > 0 {
		fromDate := time.Now().AddDate(0, 0, -days).Format("20060102T0000")
		params.Add("time_from", fromDate)
	}

	endpoint := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	slog.Info("Fetching Alpha Vantage trending news", "days", days)

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

	var apiResp alphaVantageResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	news := make([]NewsItem, 0, len(apiResp.Feed))
	for _, feed := range apiResp.Feed {
		publishedAt, err := time.Parse("20060102T150405", feed.TimePublished)
		if err != nil {
			publishedAt = time.Now()
		}

		sentiment := feed.OverallSentimentLabel
		if sentiment == "" {
			sentiment = "neutral"
		}

		news = append(news, NewsItem{
			Title:       feed.Title,
			Description: feed.Summary,
			URL:         feed.URL,
			Date:        publishedAt,
			Source:      feed.Source,
			Sentiment:   sentiment,
		})
	}

	slog.Info("Alpha Vantage trending news fetched", "count", len(news))

	return news, nil
}

func (c *AlphaVantageClient) HealthCheck() error {
	params := url.Values{}
	params.Add("function", "NEWS_SENTIMENT")
	params.Add("apikey", c.apiKey)
	params.Add("limit", "1")

	endpoint := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

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
