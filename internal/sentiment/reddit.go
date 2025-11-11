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

type RedditScraper struct {
	httpClient *http.Client
}

func NewRedditScraper() *RedditScraper {
	return &RedditScraper{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *RedditScraper) GetSentiment(ticker string) (PlatformSentiment, error) {
	subreddits := []string{"wallstreetbets", "stocks", "investing"}

	positive := 0
	negative := 0
	neutral := 0

	for _, subreddit := range subreddits {
		url := fmt.Sprintf("https://old.reddit.com/r/%s/search?q=%s&restrict_sr=on&sort=new&t=week", subreddit, ticker)

		slog.Info("Scraping Reddit", "subreddit", subreddit, "ticker", ticker)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			slog.Warn("Failed to create Reddit request", "error", err)
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

		resp, err := s.httpClient.Do(req)
		if err != nil {
			slog.Warn("Failed to fetch Reddit", "error", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			slog.Warn("Reddit returned non-200 status", "status", resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Warn("Failed to read Reddit response", "error", err)
			continue
		}

		content := string(body)

		p, n, ne := s.analyzeSentiment(content, ticker)
		positive += p
		negative += n
		neutral += ne
	}

	total := positive + negative + neutral
	if total == 0 {
		return PlatformSentiment{
			Positive: 33,
			Negative: 33,
			Neutral:  34,
		}, nil
	}

	return PlatformSentiment{
		Positive: (positive * 100) / total,
		Negative: (negative * 100) / total,
		Neutral:  (neutral * 100) / total,
	}, nil
}

func (s *RedditScraper) analyzeSentiment(content, ticker string) (positive, negative, neutral int) {
	positiveKeywords := []string{"bullish", "moon", "rocket", "buy", "calls", "long", "green", "gains", "profit", "winning"}
	negativeKeywords := []string{"bearish", "crash", "sell", "puts", "short", "red", "loss", "dump", "failing"}

	contentLower := strings.ToLower(content)
	tickerLower := strings.ToLower(ticker)

	tickerRegex := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, tickerLower))
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
