package sentiment

import (
	"log/slog"
)

type Aggregator struct {
	reddit *RedditScraper
	x      *XScraper
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		reddit: NewRedditScraper(),
		x:      NewXScraper(),
	}
}

func (a *Aggregator) GetSentiment(ticker string) (SentimentData, error) {
	redditSentiment, err := a.reddit.GetSentiment(ticker)
	if err != nil {
		slog.Warn("Reddit sentiment failed, using defaults", "error", err, "ticker", ticker)
		redditSentiment = PlatformSentiment{
			Positive: 33,
			Negative: 33,
			Neutral:  34,
		}
	}

	xSentiment, err := a.x.GetSentiment(ticker)
	if err != nil {
		slog.Warn("X sentiment failed, using defaults", "error", err, "ticker", ticker)
		xSentiment = PlatformSentiment{
			Positive: 33,
			Negative: 33,
			Neutral:  34,
		}
	}

	avgPositive := (redditSentiment.Positive + xSentiment.Positive) / 2
	avgNegative := (redditSentiment.Negative + xSentiment.Negative) / 2

	overall := "Neutral"
	if avgPositive > avgNegative+10 {
		overall = "Bullish"
	} else if avgNegative > avgPositive+10 {
		overall = "Bearish"
	}

	return SentimentData{
		Ticker:           ticker,
		XSentiment:       xSentiment,
		RedditSentiment:  redditSentiment,
		OverallSentiment: overall,
	}, nil
}
