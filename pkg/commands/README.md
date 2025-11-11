# Discord Bot Commands

This directory contains all the Go command implementations for the Discord bot.

## Available Commands

### Stock Command (`/stock`)
Retrieves latest stock news from available sources, performs X.com and Reddit sentiment analysis for user sentiment regarding certain stocks, and uses the local llama.cpp service to process data into a report of stocks to watch and why they are interesting.

#### Usage
```
/stock tickers: [ticker1,ticker2,...] [days: number_of_days]
```

#### Parameters
- `tickers` (required): Comma-separated list of stock tickers (e.g. AAPL,MSFT,GOOGL)
- `days` (optional): Number of days of news to retrieve (default: 7)

#### Features
1. Retrieves latest stock news from financial sources
2. Performs X.com (Twitter) and Reddit sentiment analysis
3. Uses llama.cpp to process data into a comprehensive report
4. Presents results in a Discord-friendly format

#### Example
```
/stock tickers: AAPL,MSFT
/stock tickers: TSLA,AMZN days: 14