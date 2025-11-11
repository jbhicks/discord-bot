#!/bin/bash

# Script to set up systemd service for Discord bot

SERVICE_FILE="/etc/systemd/system/discord-bot.service"
BOT_DIR="$(pwd)"
BOT_BINARY="$BOT_DIR/bot"
USER="$(whoami)"

# Check if running as root or with sudo
if [[ $EUID -eq 0 ]]; then
    echo "Don't run as root. Use sudo if needed."
    exit 1
fi

# Check if bot binary exists
if [[ ! -f "$BOT_BINARY" ]]; then
    echo "Bot binary not found at $BOT_BINARY. Build it first with 'go build ./cmd/bot'."
    exit 1
fi

# Create service file
sudo tee "$SERVICE_FILE" > /dev/null <<EOF
[Unit]
Description=Discord Bot
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$BOT_DIR
ExecStart=$BOT_BINARY
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

echo "Service file created at $SERVICE_FILE"

# Reload systemd
sudo systemctl daemon-reload
echo "Systemd reloaded"

# Enable and start service
sudo systemctl enable discord-bot
echo "Service enabled"

sudo systemctl start discord-bot
echo "Service started"

# Check status
sudo systemctl status discord-bot --no-pager

echo "Bot is now running as a service. Use 'sudo systemctl restart discord-bot' after updates."