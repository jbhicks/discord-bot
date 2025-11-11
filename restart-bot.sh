#!/bin/bash
sudo systemctl restart discord-bot
systemctl status discord-bot --no-pager -l | head -20
