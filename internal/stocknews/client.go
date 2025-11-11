package stocknews

import (
	"fmt"
	"log/slog"
)

type FallbackClient struct {
	primary  Client
	fallback Client
}

func NewFallbackClient(primary, fallback Client) *FallbackClient {
	return &FallbackClient{
		primary:  primary,
		fallback: fallback,
	}
}

func (c *FallbackClient) GetNews(ticker string, days int) ([]NewsItem, error) {
	news, err := c.primary.GetNews(ticker, days)
	if err == nil {
		slog.Info("Using primary news client", "ticker", ticker)
		return news, nil
	}

	slog.Warn("Primary news client failed, using fallback", "error", err, "ticker", ticker)

	news, err = c.fallback.GetNews(ticker, days)
	if err != nil {
		return nil, fmt.Errorf("both primary and fallback clients failed: %w", err)
	}

	slog.Info("Using fallback news client", "ticker", ticker)
	return news, nil
}

func (c *FallbackClient) GetTrendingNews(days int) ([]NewsItem, error) {
	news, err := c.primary.GetTrendingNews(days)
	if err == nil {
		slog.Info("Using primary news client for trending")
		return news, nil
	}

	slog.Warn("Primary news client failed for trending, using fallback", "error", err)

	news, err = c.fallback.GetTrendingNews(days)
	if err != nil {
		return nil, fmt.Errorf("both primary and fallback clients failed: %w", err)
	}

	slog.Info("Using fallback news client for trending")
	return news, nil
}

func (c *FallbackClient) HealthCheck() error {
	if err := c.primary.HealthCheck(); err != nil {
		slog.Warn("Primary client health check failed", "error", err)
		return c.fallback.HealthCheck()
	}
	return nil
}
