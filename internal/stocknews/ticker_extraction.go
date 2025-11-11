package stocknews

import (
	"regexp"
	"strings"
)

var tickerPattern = regexp.MustCompile(`\b[A-Z]{1,5}\b`)

func ExtractTickers(news []NewsItem) map[string]int {
	tickerCount := make(map[string]int)

	commonWords := map[string]bool{
		"A": true, "I": true, "THE": true, "AND": true, "OR": true, "BUT": true,
		"IN": true, "ON": true, "AT": true, "TO": true, "FOR": true, "OF": true,
		"AS": true, "BY": true, "AN": true, "BE": true, "IS": true, "IT": true,
		"CEO": true, "CFO": true, "CTO": true, "IPO": true, "SEC": true, "USA": true,
		"US": true, "UK": true, "EU": true, "AI": true, "API": true, "APP": true,
		"NEW": true, "ALL": true, "GET": true, "SET": true, "HAS": true, "HAD": true,
		"WAS": true, "ARE": true, "CAN": true, "MAY": true, "NOT": true, "OUT": true,
		"UP": true, "SO": true, "NO": true, "IF": true, "MY": true, "ONE": true,
		"TWO": true, "TOP": true, "BIG": true, "OLD": true, "NOW": true, "OUR": true,
		"ITS": true, "KEY": true, "WHY": true, "HOW": true, "WHO": true, "GOT": true,
		"ETF": true, "FBI": true, "CIA": true, "DOJ": true, "FDA": true,
		"WITH": true, "FROM": true, "INTO": true, "OVER": true, "THAN": true,
		"THIS": true, "THAT": true, "HAVE": true, "WILL": true, "BEEN": true,
		"WERE": true, "THEY": true, "THEM": true, "YOUR": true, "THEIR": true,
		"ORDER": true, "FIRST": true, "LAST": true, "NEXT": true, "YEAR": true,
		"TIME": true, "WEEK": true, "MONTH": true, "PRICE": true, "STOCK": true,
	}

	for _, item := range news {
		text := item.Title + " " + item.Description
		text = strings.ToUpper(text)

		matches := tickerPattern.FindAllString(text, -1)
		for _, match := range matches {
			if len(match) < 2 || len(match) > 5 {
				continue
			}
			if commonWords[match] {
				continue
			}
			tickerCount[match]++
		}
	}

	return tickerCount
}

func GetTopTickers(tickerCount map[string]int, limit int) []string {
	type tickerFreq struct {
		ticker string
		count  int
	}

	freqs := make([]tickerFreq, 0, len(tickerCount))
	for ticker, count := range tickerCount {
		freqs = append(freqs, tickerFreq{ticker, count})
	}

	for i := 0; i < len(freqs); i++ {
		for j := i + 1; j < len(freqs); j++ {
			if freqs[j].count > freqs[i].count {
				freqs[i], freqs[j] = freqs[j], freqs[i]
			}
		}
	}

	result := make([]string, 0, limit)
	for i := 0; i < len(freqs) && i < limit; i++ {
		result = append(result, freqs[i].ticker)
	}

	return result
}
