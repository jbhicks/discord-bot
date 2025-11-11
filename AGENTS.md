# Agent Guidelines for Go Discord Bot Project

## Commands
- **Build**: `go build -o bin/bot ./cmd/bot`
- **Dev**: `go run ./cmd/bot` (or use air for hot reloading)
- **Start**: `./bin/bot`
- **Test**: `go test ./...`
- **Lint/Format**: `gofmt -w .` and `go vet ./...`
- **Mod Tidy**: `go mod tidy`

## Code Style
- **Modules**: Go modules with `go.mod` (use `github.com/username/discord-bot`)
- **Commands**: Implement `Command` interface with `Name()`, `Description()`, `Execute()` methods
- **Naming**: PascalCase for exported types/functions, camelCase for unexported
- **Imports**: Standard library first, then third-party, then local packages
- **Error Handling**: Return errors, use `if err != nil` pattern, log with structured logging
- **Logging**: Use `slog` or `log/slog` for structured logging, avoid `fmt.Printf`
  - Log all command invocations with user, user_id, guild_id, and relevant parameters
  - Log before deferred responses to capture user intent early
- **Types**: Define interfaces for contracts, use structs for data, leverage Go's type system
- **Formatting**: `gofmt` (no configuration needed), tabs for indentation
- **Comments**: Go doc comments for exported functions/types, minimal inline comments
- **Concurrency**: Use goroutines and channels for concurrent operations
- **AI Integration**: 
  - Use llama.cpp server for local LLM responses (port 8081)
  - Use Stable Diffusion WebUI for image generation (port 7860)
  - Handle in separate goroutines for long-running operations

## Available Commands
- `/ping` - Check bot latency
- `/play` - Play audio in voice channel
- `/stop` - Stop playback
- `/queue` - Show playback queue
- `/skip` - Skip current track
- `/loop` - Toggle loop mode
- `/search` - Search for content
- `/playlist` - Manage playlists
- `/ai` - Ask AI a question (uses llama.cpp)
- `/imagine` - Generate images from text prompts (uses Stable Diffusion)
- `/help` - Show help information

## Services
- **Discord Bot**: Port varies, systemd service `discord-bot`
- **Llama.cpp Server**: Port 8081, systemd service `llama-server`
- **SD WebUI**: Port 7860, systemd service `stable-diffusion-webui`

## Important Notes
- **NEVER use `curl -f` without timeout**: Always use `curl -f --max-time 5` or `curl -f -m 5` to prevent blocking
- When testing web services, always set a timeout to avoid infinite hangs