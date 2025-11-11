package commands

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/josh/discord-bot/internal/llm"
	"github.com/josh/discord-bot/internal/sentiment"
	"github.com/josh/discord-bot/internal/stocknews"
)

type LLMClient interface {
	Chat(prompt string) (string, error)
}

type SentimentClient interface {
	GetSentiment(ticker string) (sentiment.SentimentData, error)
}

type StockCommand struct {
	newsClient      stocknews.Client
	llmClient       LLMClient
	sentimentClient SentimentClient
}

func NewStockCommand(newsClient stocknews.Client, llmClient *llm.Client, sentimentClient *sentiment.Aggregator) *StockCommand {
	return &StockCommand{
		newsClient:      newsClient,
		llmClient:       llmClient,
		sentimentClient: sentimentClient,
	}
}

// Name returns the name of the command
func (c *StockCommand) Name() string {
	return "stock"
}

// Description returns the description of the command
func (c *StockCommand) Description() string {
	return "Get stock news, sentiment analysis, and AI-generated recommendations"
}

// Data returns the command data for Discord
func (c *StockCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.Name(),
		Description: c.Description(),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "tickers",
				Description: "Comma-separated list of stock tickers (e.g. AAPL,MSFT,GOOGL). Leave empty for trending.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "days",
				Description: "Number of days of news to retrieve (default: 7)",
				Required:    false,
			},
		},
	}
}

// Execute executes the stock command
func (c *StockCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options

	var tickers string
	days := int64(7)

	for _, opt := range options {
		switch opt.Name {
		case "tickers":
			tickers = opt.StringValue()
		case "days":
			days = opt.IntValue()
		}
	}

	var tickerList []string
	if tickers != "" {
		tickerList = strings.Split(strings.TrimSpace(tickers), ",")
		for i, ticker := range tickerList {
			tickerList[i] = strings.ToUpper(strings.TrimSpace(ticker))
		}
	}

	username := "Unknown"
	userID := "Unknown"
	if i.Member != nil && i.Member.User != nil {
		username = i.Member.User.Username
		userID = i.Member.User.ID
	}

	slog.Info("Stock command received",
		"user", username,
		"user_id", userID,
		"guild_id", i.GuildID,
		"tickers", strings.Join(tickerList, ","),
		"days", days,
	)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		slog.Error("Error responding to interaction", "error", err)
		return err
	}

	go func() {
		var report string
		var err error

		if len(tickerList) == 0 {
			report, err = c.processTrendingAnalysis(int(days))
		} else {
			report, err = c.processStockAnalysis(tickerList, int(days))
		}

		if err != nil {
			slog.Error("Error processing stock analysis", "error", err)
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Content: fmt.Sprintf("Error processing stock analysis: %v", err),
			})
			return
		}

		if len(report) > 2000 {
			chunks := make([]string, 0)
			for len(report) > 0 {
				end := 2000
				if end > len(report) {
					end = len(report)
				}
				chunks = append(chunks, report[:end])
				report = report[end:]
			}

			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &chunks[0],
			})

			for _, chunk := range chunks[1:] {
				s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
					Content: chunk,
				})
			}
		} else {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &report,
			})
		}
	}()

	return nil
}

// processStockAnalysis processes the stock analysis for given tickers
func (c *StockCommand) processStockAnalysis(tickers []string, days int) (string, error) {
	var report strings.Builder

	report.WriteString("## 📊 Stock Analysis Report\n\n")

	for i, ticker := range tickers {
		if i > 0 {
			report.WriteString("\n---\n\n")
		}

		report.WriteString(fmt.Sprintf("### %s\n\n", ticker))

		// Get news for this ticker
		news, err := c.getStockNews(ticker, days)
		if err != nil {
			slog.Error("Error getting news", "ticker", ticker, "error", err)
			report.WriteString(fmt.Sprintf("Error retrieving news for %s: %v\n\n", ticker, err))
			continue
		}

		// Get sentiment analysis
		sentiment, err := c.getSentimentAnalysis(ticker)
		if err != nil {
			slog.Error("Error getting sentiment", "ticker", ticker, "error", err)
			report.WriteString(fmt.Sprintf("Error retrieving sentiment for %s: %v\n\n", ticker, err))
			continue
		}

		// Combine data and generate AI report
		aiPrompt := c.generateAIPrompt(ticker, news, sentiment)
		aiReport, err := c.generateAIReport(aiPrompt)
		if err != nil {
			slog.Error("Error generating AI report", "ticker", ticker, "error", err)
			report.WriteString(fmt.Sprintf("Error generating AI report for %s: %v\n\n", ticker, err))
			continue
		}

		report.WriteString(aiReport)
	}

	return report.String(), nil
}

// getStockNews retrieves news for a given stock ticker
func (c *StockCommand) getStockNews(ticker string, days int) ([]stocknews.NewsItem, error) {
	return c.newsClient.GetNews(ticker, days)
}

// getSentimentAnalysis retrieves sentiment analysis for a given stock ticker
func (c *StockCommand) getSentimentAnalysis(ticker string) (sentiment.SentimentData, error) {
	return c.sentimentClient.GetSentiment(ticker)
}

// generateAIPrompt generates a prompt for the AI to process the stock data
func (c *StockCommand) generateAIPrompt(ticker string, news []stocknews.NewsItem, sentimentData sentiment.SentimentData) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Analyze the following stock data for %s and provide a comprehensive report:\n\n", ticker))

	prompt.WriteString("## News Articles:\n")
	for _, article := range news {
		prompt.WriteString(fmt.Sprintf("- %s: %s\n", article.Title, article.Description))
	}

	prompt.WriteString("\n## Sentiment Analysis:\n")
	prompt.WriteString(fmt.Sprintf("X (Twitter) Sentiment - Positive: %d%%, Negative: %d%%, Neutral: %d%%\n",
		sentimentData.XSentiment.Positive, sentimentData.XSentiment.Negative, sentimentData.XSentiment.Neutral))
	prompt.WriteString(fmt.Sprintf("Reddit Sentiment - Positive: %d%%, Negative: %d%%, Neutral: %d%%\n",
		sentimentData.RedditSentiment.Positive, sentimentData.RedditSentiment.Negative, sentimentData.RedditSentiment.Neutral))
	prompt.WriteString(fmt.Sprintf("Overall Sentiment: %s\n", sentimentData.OverallSentiment))

	prompt.WriteString("\n\nProvide a concise analysis of whether this stock is worth watching, including:\n")
	prompt.WriteString("1. Key factors from news\n")
	prompt.WriteString("2. Sentiment summary\n")
	prompt.WriteString("3. Investment recommendation\n")
	prompt.WriteString("4. Risk factors\n")

	return prompt.String()
}

// generateAIReport generates a report using the llama.cpp service
func (c *StockCommand) generateAIReport(prompt string) (string, error) {
	return c.llmClient.Chat(prompt)
}

func (c *StockCommand) processTrendingAnalysis(days int) (string, error) {
	var report strings.Builder

	report.WriteString("## 📈 Trending Stocks Analysis\n\n")

	news, err := c.newsClient.GetTrendingNews(days)
	if err != nil {
		return "", fmt.Errorf("failed to get trending news: %w", err)
	}

	tickerCounts := stocknews.ExtractTickers(news)
	topTickers := stocknews.GetTopTickers(tickerCounts, 5)

	if len(topTickers) == 0 {
		return "No trending stocks found in recent news.", nil
	}

	report.WriteString("**Most mentioned tickers in the last ")
	report.WriteString(fmt.Sprintf("%d days:**\n", days))
	for i, ticker := range topTickers {
		report.WriteString(fmt.Sprintf("%d. %s (%d mentions)\n", i+1, ticker, tickerCounts[ticker]))
	}
	report.WriteString("\n---\n\n")

	for i, ticker := range topTickers {
		if i > 0 {
			report.WriteString("\n---\n\n")
		}

		report.WriteString(fmt.Sprintf("### %s\n\n", ticker))

		tickerNews, err := c.getStockNews(ticker, days)
		if err != nil {
			slog.Error("Error getting news", "ticker", ticker, "error", err)
			report.WriteString(fmt.Sprintf("Error retrieving news for %s: %v\n\n", ticker, err))
			continue
		}

		sentimentData, err := c.getSentimentAnalysis(ticker)
		if err != nil {
			slog.Error("Error getting sentiment", "ticker", ticker, "error", err)
			report.WriteString(fmt.Sprintf("Error retrieving sentiment for %s: %v\n\n", ticker, err))
			continue
		}

		aiPrompt := c.generateAIPrompt(ticker, tickerNews, sentimentData)
		aiReport, err := c.generateAIReport(aiPrompt)
		if err != nil {
			slog.Error("Error generating AI report", "ticker", ticker, "error", err)
			report.WriteString(fmt.Sprintf("Error generating AI report for %s: %v\n\n", ticker, err))
			continue
		}

		report.WriteString(aiReport)
	}

	return report.String(), nil
}
