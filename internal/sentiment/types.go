package sentiment

type SentimentData struct {
	Ticker           string            `json:"ticker"`
	XSentiment       PlatformSentiment `json:"x_sentiment"`
	RedditSentiment  PlatformSentiment `json:"reddit_sentiment"`
	OverallSentiment string            `json:"overall_sentiment"`
}

type PlatformSentiment struct {
	Positive int `json:"positive"`
	Negative int `json:"negative"`
	Neutral  int `json:"neutral"`
}
