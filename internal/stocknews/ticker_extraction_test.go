package stocknews

import (
	"testing"
	"time"
)

func TestExtractTickers(t *testing.T) {
	news := []NewsItem{
		{
			Title:       "AAPL Reports Strong Q4 Earnings",
			Description: "Apple (AAPL) beat expectations with revenue driven by iPhone sales.",
			Date:        time.Now(),
		},
		{
			Title:       "MSFT Azure Growth Continues",
			Description: "Microsoft (MSFT) announces strong cloud growth.",
			Date:        time.Now(),
		},
		{
			Title:       "AAPL and GOOGL Partner on AI",
			Description: "Apple and Google announce AI partnership.",
			Date:        time.Now(),
		},
	}

	tickers := ExtractTickers(news)

	if tickers["AAPL"] != 3 {
		t.Errorf("Expected AAPL to appear 3 times, got %d", tickers["AAPL"])
	}

	if tickers["MSFT"] != 2 {
		t.Errorf("Expected MSFT to appear 2 times, got %d", tickers["MSFT"])
	}

	if tickers["GOOGL"] != 1 {
		t.Errorf("Expected GOOGL to appear 1 time, got %d", tickers["GOOGL"])
	}

	if tickers["THE"] > 0 {
		t.Errorf("Common word 'THE' should be filtered out")
	}
}

func TestGetTopTickers(t *testing.T) {
	tickerCount := map[string]int{
		"AAPL":  10,
		"MSFT":  8,
		"GOOGL": 6,
		"TSLA":  4,
		"AMZN":  2,
	}

	top3 := GetTopTickers(tickerCount, 3)

	if len(top3) != 3 {
		t.Fatalf("Expected 3 tickers, got %d", len(top3))
	}

	if top3[0] != "AAPL" {
		t.Errorf("Expected AAPL to be first, got %s", top3[0])
	}

	if top3[1] != "MSFT" {
		t.Errorf("Expected MSFT to be second, got %s", top3[1])
	}

	if top3[2] != "GOOGL" {
		t.Errorf("Expected GOOGL to be third, got %s", top3[2])
	}
}

func TestGetTopTickersMoreThanAvailable(t *testing.T) {
	tickerCount := map[string]int{
		"AAPL": 10,
		"MSFT": 8,
	}

	top5 := GetTopTickers(tickerCount, 5)

	if len(top5) != 2 {
		t.Fatalf("Expected 2 tickers when requesting 5 with only 2 available, got %d", len(top5))
	}
}
