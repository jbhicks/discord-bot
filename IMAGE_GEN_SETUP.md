# Discord Bot - Image Generation Setup Guide

This guide will help you set up the image generation feature for the Discord bot using Stable Diffusion WebUI.

## Prerequisites

- AMD Strix Halo with 128GB RAM (or other system with sufficient resources)
- Python 3.10 or 3.11
- Git
- Sufficient disk space (~20GB for WebUI + models)

## Installation Steps

### 1. Run the Setup Script

The repository includes a setup script that will install Stable Diffusion WebUI and download the SDXL-Turbo model:

```bash
cd ~/discord-bot
./setup-sd-webui.sh
```

This will:
- Clone the Stable Diffusion WebUI repository to `~/stable-diffusion-webui`
- Download the SDXL-Turbo model (~7GB)
- Prepare the environment

### 2. First-Time WebUI Setup

Run the WebUI once manually to complete the installation:

```bash
cd ~/stable-diffusion-webui
./webui.sh --api --listen --port 7860
```

This will:
- Install Python dependencies
- Set up the virtual environment
- Download additional required files

Wait for it to complete and show "Running on local URL: http://0.0.0.0:7860"

Press Ctrl+C to stop it once ready.

### 3. Configure Environment Variable

Add the image generation URL to your `.env` file (optional, defaults to http://localhost:7860):

```bash
echo "IMAGE_GEN_URL=http://localhost:7860" >> ~/discord-bot/.env
```

### 4. Install as System Service

To run the image generation server automatically:

```bash
sudo cp ~/discord-bot/sd-webui.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable sd-webui
sudo systemctl start sd-webui
```

Check the status:

```bash
sudo systemctl status sd-webui
```

View logs:

```bash
sudo journalctl -u sd-webui -f
```

### 5. Rebuild and Deploy the Bot

```bash
cd ~/discord-bot
go build -o bin/bot ./cmd/bot
```

If running as a service:

```bash
sudo systemctl restart discord-bot
```

Or run manually:

```bash
./bin/bot
```

## Using the /imagine Command

Once everything is running, you can use the new `/imagine` command in Discord:

### Basic Usage

```
/imagine prompt: a beautiful sunset over mountains
```

### Advanced Options

```
/imagine prompt: a cyberpunk city at night
         width: 768
         height: 512
         steps: 8
```

**Parameters:**
- `prompt` (required): Description of the image to generate
- `width` (optional): Image width in pixels (default: 512)
- `height` (optional): Image height in pixels (default: 512)
- `steps` (optional): Number of generation steps (default: 4 for turbo)

**Note:** With SDXL-Turbo, 4 steps is optimal. More steps may actually reduce quality.

## Performance

With your AMD Strix Halo (128GB RAM):

- **SDXL-Turbo**: ~2-5 seconds per image (512x512, 4 steps)
- **SDXL-Turbo**: ~5-10 seconds per image (768x768, 4 steps)
- CPU inference is viable with your large RAM

## Upgrading to Better Models

### FLUX.1-schnell (Faster, Better Quality)

1. Download the model:
```bash
cd ~/stable-diffusion-webui/models/Stable-diffusion
wget https://huggingface.co/black-forest-labs/FLUX.1-schnell/resolve/main/flux1-schnell.safetensors
```

2. Restart the service:
```bash
sudo systemctl restart sd-webui
```

3. Select the model in the WebUI or via API

### FLUX.1-dev (Best Quality)

For the highest quality (slower):

```bash
cd ~/stable-diffusion-webui/models/Stable-diffusion
wget https://huggingface.co/black-forest-labs/FLUX.1-dev/resolve/main/flux1-dev.safetensors
```

Adjust steps to 20-50 for FLUX.1-dev in the command.

## Troubleshooting

### Image Generation Server Not Responding

Check if the service is running:
```bash
sudo systemctl status sd-webui
```

Test the API manually:
```bash
curl http://localhost:7860/sdapi/v1/sd-models
```

### Out of Memory Errors

If you get OOM errors, try:
- Reducing image size (e.g., 512x512 instead of 1024x1024)
- Using SDXL-Turbo instead of larger models
- Adjusting the `--no-half` flag in the service file

### Slow Generation

- SDXL-Turbo should be fast (2-10 seconds)
- Check CPU usage with `htop`
- Ensure no other heavy processes are running
- Consider GPU acceleration if available (ROCm for AMD)

## Architecture

```
Discord User
    ↓
Discord Bot (Go)
    ↓ HTTP Request
Image Generation Client (internal/imagegen)
    ↓ POST /sdapi/v1/txt2img
Stable Diffusion WebUI (Python)
    ↓
SDXL-Turbo Model
    ↓
Generated Image (Base64)
    ↓ Upload
Discord (Image Display)
```

## Logging

The bot logs all image generation requests with:
- Username and User ID
- Guild ID
- Prompt
- Generation parameters
- Success/failure status

Check bot logs for debugging:
```bash
sudo journalctl -u discord-bot -f
```

## Security Considerations

- The image generation server runs locally (localhost only by default)
- Consider adding prompt filtering for inappropriate content
- Monitor logs for misuse
- Rate limiting is recommended for production use

## Next Steps

Future enhancements could include:
- Image-to-image transformation (`/reimagine`)
- Style presets (`/style`)
- Model switching
- Queue system for concurrent requests
- Upscaling capabilities
- Negative prompt customization

Enjoy generating images! 🎨
