.PHONY: build dev start test lint fmt clean logs

build:
	go build -o bin/bot ./cmd/bot
	sudo systemctl restart discord-bot.service

dev:
	air

start:
	./bin/bot

test:
	go test ./...

lint:
	go vet ./...

fmt:
	gofmt -w .

clean:
	rm -f bin/bot

logs:
	sudo journalctl -u discord-bot.service -f
