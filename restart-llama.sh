#!/bin/bash

sudo pkill -9 llama-server
sleep 2
sudo rm -f /tmp/llama-server.log
nohup /home/josh/llama.cpp/build-vk/bin/llama-server -m /home/josh/models/Huihui-Qwen3-Coder-30B-A3B-Instruct-abliterated.Q4_K_M.gguf --host 0.0.0.0 --port 8081 -c 32768 -ngl 999 -fa on -b 256 -ub 4096 --jinja > /tmp/llama-server.log 2>&1 &
sleep 5
curl -s http://localhost:8081/v1/models | python3 -c "import sys, json; d=json.load(sys.stdin); print(d['data'][0]['id'])"
