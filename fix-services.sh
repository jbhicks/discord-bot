#!/bin/bash

# Copy updated service file
cp llama-server.service /etc/systemd/system/

# Reload systemd
systemctl daemon-reload

# Stop and start the server service to ensure changes take effect
systemctl stop llama-server.service
systemctl start llama-server.service

echo "Services updated and restarted. Check with: checkservices"