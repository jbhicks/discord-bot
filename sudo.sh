#!/bin/bash

cp llama-server.service /etc/systemd/system/
cp discord-bot.service /etc/systemd/system/
cp sd-webui.service /etc/systemd/system/

systemctl daemon-reload

systemctl enable llama-server.service
systemctl enable discord-bot.service
systemctl enable sd-webui.service

systemctl start llama-server.service
systemctl start discord-bot.service
systemctl start sd-webui.service

echo "Services enabled and started. Check status with:"
echo "systemctl status llama-server.service"
echo "systemctl status discord-bot.service"
echo "systemctl status sd-webui.service"