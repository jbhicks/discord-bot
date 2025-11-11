package commands

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/josh/discord-bot/internal/sentiment"
	"github.com/josh/discord-bot/internal/stocknews"
)

type mockNewsClient struct{}

func (m *mockNewsClient) GetNews(ticker string, days int) ([]stocknews.NewsItem, error) {
	return []stocknews.NewsItem{
		{
			Title:       fmt.Sprintf("%s Beats Earnings Expectations", ticker),
			Description: fmt.Sprintf("%s reported Q4 earnings that exceeded analyst expectations by 15%%, driven by strong product sales and improved margins.", ticker),
			URL:         fmt.Sprintf("https://example.com/news/%s-earnings", ticker),
			Date:        time.Now().AddDate(0, 0, -1),
			Source:      "Financial Times",
			Sentiment:   "positive",
		},
		{
			Title:       fmt.Sprintf("%s Announces Strategic Partnership", ticker),
			Description: fmt.Sprintf("%s announced a major partnership that is expected to drive 20%% revenue growth in the next fiscal year.", ticker),
			URL:         fmt.Sprintf("https://example.com/news/%s-partnership", ticker),
			Date:        time.Now().AddDate(0, 0, -2),
			Source:      "Bloomberg",
			Sentiment:   "positive",
		},
		{
			Title:       fmt.Sprintf("%s Faces Regulatory Scrutiny", ticker),
			Description: "Regulatory concerns raised about the company's business practices in certain markets.",
			URL:         fmt.Sprintf("https://example.com/news/%s-regulatory", ticker),
			Date:        time.Now().AddDate(0, 0, -3),
			Source:      "Reuters",
			Sentiment:   "negative",
		},
	}, nil
}

func (m *mockNewsClient) GetTrendingNews(days int) ([]stocknews.NewsItem, error) {
	return []stocknews.NewsItem{
		{
			Title:       "AAPL Launches New Product Line",
			Description: "Apple (AAPL) announces revolutionary new products with strong pre-order numbers.",
			URL:         "https://example.com/news/aapl-product",
			Date:        time.Now(),
			Source:      "TechCrunch",
			Sentiment:   "positive",
		},
		{
			Title:       "AAPL Stock Surges",
			Description: "AAPL shares hit all-time high on strong earnings report.",
			URL:         "https://example.com/news/aapl-surge",
			Date:        time.Now(),
			Source:      "Bloomberg",
			Sentiment:   "positive",
		},
		{
			Title:       "MSFT Cloud Growth Accelerates",
			Description: "Microsoft (MSFT) Azure sees 50% growth in enterprise adoption.",
			URL:         "https://example.com/news/msft-cloud",
			Date:        time.Now(),
			Source:      "Forbes",
			Sentiment:   "positive",
		},
		{
			Title:       "GOOGL AI Investment",
			Description: "Google (GOOGL) announces major AI infrastructure investment.",
			URL:         "https://example.com/news/googl-ai",
			Date:        time.Now(),
			Source:      "Reuters",
			Sentiment:   "positive",
		},
	}, nil
}

func (m *mockNewsClient) HealthCheck() error {
	return nil
}

type mockLLMClient struct{}

func (m *mockLLMClient) Chat(prompt string) (string, error) {
	return `### Investment Analysis

**Key Findings:**
- Strong recent earnings performance with 15% beat on expectations
- Positive momentum from strategic partnership announcement
- Some regulatory headwinds that could impact short-term performance

**Sentiment Summary:**
- Social media sentiment is predominantly bullish (68% positive)
- Reddit communities showing strong retail investor interest
- X.com discussions indicate growing institutional attention

**Investment Recommendation:**
This stock shows strong fundamentals with positive earnings surprises and strategic growth initiatives. The bullish sentiment from both retail and institutional investors suggests continued upward momentum. However, investors should monitor regulatory developments closely.

**Risk Factors:**
- Regulatory scrutiny could lead to fines or business restrictions
- Market volatility may affect short-term price action
- Valuation metrics suggest the stock may be trading at premium levels

**Conclusion:**
Based on the analysis, this appears to be a solid investment opportunity with good growth potential, but investors should maintain appropriate position sizing given the regulatory risks.`, nil
}

type mockSentimentAggregator struct{}

func (m *mockSentimentAggregator) GetSentiment(ticker string) (sentiment.SentimentData, error) {
	return sentiment.SentimentData{
		Ticker: ticker,
		XSentiment: sentiment.PlatformSentiment{
			Positive: 65,
			Negative: 20,
			Neutral:  15,
		},
		RedditSentiment: sentiment.PlatformSentiment{
			Positive: 72,
			Negative: 15,
			Neutral:  13,
		},
		OverallSentiment: "Bullish",
	}, nil
}

func TestStockCommand_ProcessAnalysis(t *testing.T) {
	newsClient := &mockNewsClient{}
	llmClient := &mockLLMClient{}
	sentimentClient := &mockSentimentAggregator{}

	cmd := &StockCommand{
		newsClient:      newsClient,
		llmClient:       llmClient,
		sentimentClient: sentimentClient,
	}

	tickers := []string{"AAPL", "MSFT"}
	days := 7

	t.Logf("\n=== Testing Stock Analysis for: %v (last %d days) ===\n", tickers, days)

	report, err := cmd.processStockAnalysis(tickers, days)
	if err != nil {
		t.Fatalf("processStockAnalysis failed: %v", err)
	}

	if report == "" {
		t.Error("Report should not be empty")
	}

	t.Logf("\n=== GENERATED REPORT ===\n\n%s\n\n=== END REPORT ===\n", report)

	if len(report) > 2000 {
		chunks := (len(report) + 1999) / 2000
		t.Logf("\nNote: Report is %d characters, would be split into %d Discord messages", len(report), chunks)
	} else {
		t.Logf("\nReport length: %d characters (fits in single Discord message)", len(report))
	}
}

func TestStockCommand_SingleTicker(t *testing.T) {
	newsClient := &mockNewsClient{}
	llmClient := &mockLLMClient{}
	sentimentClient := &mockSentimentAggregator{}

	cmd := &StockCommand{
		newsClient:      newsClient,
		llmClient:       llmClient,
		sentimentClient: sentimentClient,
	}

	ticker := "AAPL"
	days := 7

	t.Logf("\n=== Testing Single Ticker Analysis: %s ===\n", ticker)

	report, err := cmd.processStockAnalysis([]string{ticker}, days)
	if err != nil {
		t.Fatalf("processStockAnalysis failed: %v", err)
	}

	t.Logf("\n%s\n", report)
}

func TestStockCommand_TrendingAnalysis(t *testing.T) {
	newsClient := &mockNewsClient{}
	llmClient := &mockLLMClient{}
	sentimentClient := &mockSentimentAggregator{}

	cmd := &StockCommand{
		newsClient:      newsClient,
		llmClient:       llmClient,
		sentimentClient: sentimentClient,
	}

	days := 7

	t.Logf("\n=== Testing Trending Analysis (no tickers provided) ===\n")

	report, err := cmd.processTrendingAnalysis(days)
	if err != nil {
		t.Fatalf("processTrendingAnalysis failed: %v", err)
	}

	if report == "" {
		t.Error("Report should not be empty")
	}

	if !strings.Contains(report, "Trending Stocks Analysis") {
		t.Error("Report should contain 'Trending Stocks Analysis' header")
	}

	if !strings.Contains(report, "AAPL") {
		t.Error("Report should contain AAPL ticker from trending news")
	}

	if !strings.Contains(report, "MSFT") {
		t.Error("Report should contain MSFT ticker from trending news")
	}

	if !strings.Contains(report, "mentions") {
		t.Error("Report should show mention counts for trending tickers")
	}

	t.Logf("\n=== GENERATED TRENDING REPORT ===\n\n%s\n\n=== END REPORT ===\n", report)

	if len(report) > 2000 {
		chunks := (len(report) + 1999) / 2000
		t.Logf("\nNote: Report is %d characters, would be split into %d Discord messages", len(report), chunks)
	} else {
		t.Logf("\nReport length: %d characters (fits in single Discord message)", len(report))
	}
}
