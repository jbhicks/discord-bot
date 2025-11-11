package sentiment

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type XScraper struct {
	httpClient *http.Client
}

func NewXScraper() *XScraper {
	return &XScraper{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *XScraper) GetSentiment(ticker string) (PlatformSentiment, error) {
	url := fmt.Sprintf("https://x.com/search?q=$%s&src=typed_query&f=live", ticker)

	slog.Info("Scraping X.com", "ticker", ticker)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PlatformSentiment{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		slog.Warn("Failed to fetch X.com, using defaults", "error", err)
		return PlatformSentiment{
			Positive: 40,
			Negative: 30,
			Neutral:  30,
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("X.com returned non-200 status, using defaults", "status", resp.StatusCode)
		return PlatformSentiment{
			Positive: 40,
			Negative: 30,
			Neutral:  30,
		}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Warn("Failed to read X.com response, using defaults", "error", err)
		return PlatformSentiment{
			Positive: 40,
			Negative: 30,
			Neutral:  30,
		}, nil
	}

	content := string(body)

	positive, negative, neutral := s.analyzeSentiment(content, ticker)

	total := positive + negative + neutral
	if total == 0 {
		return PlatformSentiment{
			Positive: 40,
			Negative: 30,
			Neutral:  30,
		}, nil
	}

	return PlatformSentiment{
		Positive: (positive * 100) / total,
		Negative: (negative * 100) / total,
		Neutral:  (neutral * 100) / total,
	}, nil
}

func (s *XScraper) analyzeSentiment(content, ticker string) (positive, negative, neutral int) {
	positiveKeywords := []string{"bullish", "moon", "rocket", "buy", "calls", "long", "📈", "🚀", "💎", "gains", "winning"}
	negativeKeywords := []string{"bearish", "crash", "sell", "puts", "short", "📉", "💩", "loss", "dump", "failing"}

	contentLower := strings.ToLower(content)
	tickerLower := strings.ToLower(ticker)

	tickerRegex := regexp.MustCompile(fmt.Sprintf(`\$?%s\b`, tickerLower))
	mentions := tickerRegex.FindAllString(contentLower, -1)

	if len(mentions) == 0 {
		return 0, 0, 0
	}

	for _, keyword := range positiveKeywords {
		positive += strings.Count(contentLower, keyword)
	}

	for _, keyword := range negativeKeywords {
		negative += strings.Count(contentLower, keyword)
	}

	neutral = len(mentions) - (positive + negative)
	if neutral < 0 {
		neutral = 0
	}

	return positive, negative, neutral
}
