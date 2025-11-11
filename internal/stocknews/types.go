package stocknews

import "time"

type NewsItem struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Date        time.Time `json:"date"`
	Source      string    `json:"source"`
	Sentiment   string    `json:"sentiment"`
}

type Client interface {
	GetNews(ticker string, days int) ([]NewsItem, error)
	GetTrendingNews(days int) ([]NewsItem, error)
	HealthCheck() error
}
