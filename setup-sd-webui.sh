#!/bin/bash
set -e

echo "Setting up Stable Diffusion WebUI for Discord Bot Image Generation"
echo "===================================================================="
echo ""

SD_DIR="$HOME/stable-diffusion-webui"
MODEL_DIR="$SD_DIR/models/Stable-diffusion"

# Clone the repository
if [ ! -d "$SD_DIR" ]; then
    echo "Cloning Stable Diffusion WebUI..."
    cd "$HOME"
    git clone https://github.com/AUTOMATIC1111/stable-diffusion-webui.git
else
    echo "Stable Diffusion WebUI already exists at $SD_DIR"
fi

# Download SDXL-Turbo model
echo ""
echo "Downloading SDXL-Turbo model..."
echo "This is a ~7GB download and may take a while..."
echo ""

mkdir -p "$MODEL_DIR"

if [ ! -f "$MODEL_DIR/sd_xl_turbo_1.0_fp16.safetensors" ]; then
    cd "$MODEL_DIR"
    wget -O sd_xl_turbo_1.0_fp16.safetensors \
        "https://huggingface.co/stabilityai/sdxl-turbo/resolve/main/sd_xl_turbo_1.0_fp16.safetensors"
    echo "Model downloaded successfully!"
else
    echo "Model already exists, skipping download"
fi

echo ""
echo "Setup complete! Next steps:"
echo "1. Run the WebUI once to complete installation:"
echo "   cd $SD_DIR && ./webui.sh --api --listen --port 7860"
echo ""
echo "2. After testing, create systemd service:"
echo "   sudo cp ~/discord-bot/sd-webui.service /etc/systemd/system/"
echo "   sudo systemctl daemon-reload"
echo "   sudo systemctl enable sd-webui"
echo "   sudo systemctl start sd-webui"
echo ""
