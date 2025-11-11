#!/bin/bash
set -e

echo "================================================"
echo "SDXL Base 1.0 Installation Script"
echo "================================================"
echo ""

MODEL_DIR="/home/josh/stable-diffusion-webui/models/Stable-diffusion"
MODEL_NAME="sd_xl_base_1.0.safetensors"
MODEL_PATH="$MODEL_DIR/$MODEL_NAME"
HUGGINGFACE_REPO="stabilityai/stable-diffusion-xl-base-1.0"
HUGGINGFACE_FILE="sd_xl_base_1.0.safetensors"

# Check if model already exists
if [ -f "$MODEL_PATH" ]; then
    echo "✓ SDXL Base 1.0 already downloaded at $MODEL_PATH"
else
    echo "Downloading SDXL Base 1.0 (6.9GB)..."
    echo "This may take several minutes depending on your connection."
    echo ""
    
    # Download using huggingface-cli
    cd "$MODEL_DIR"
    huggingface-cli download "$HUGGINGFACE_REPO" "$HUGGINGFACE_FILE" --local-dir . --local-dir-use-symlinks False
    
    echo ""
    echo "✓ Download complete!"
fi

echo ""
echo "Switching model via SD WebUI API..."

# Get current model
CURRENT_MODEL=$(curl -s http://localhost:7860/sdapi/v1/options | python3 -c "import sys, json; print(json.load(sys.stdin).get('sd_model_checkpoint', 'unknown'))")
echo "Current model: $CURRENT_MODEL"

# Switch to SDXL Base 1.0
curl -s -X POST http://localhost:7860/sdapi/v1/options \
  -H "Content-Type: application/json" \
  -d "{\"sd_model_checkpoint\": \"$MODEL_NAME\"}" > /dev/null

echo "Waiting for model to load..."
sleep 5

# Verify switch
NEW_MODEL=$(curl -s http://localhost:7860/sdapi/v1/options | python3 -c "import sys, json; print(json.load(sys.stdin).get('sd_model_checkpoint', 'unknown'))")
echo "New model: $NEW_MODEL"

echo ""
echo "================================================"
echo "✓ Successfully switched to SDXL Base 1.0"
echo "================================================"
echo ""
echo "Model details:"
echo "  - Better quality than Turbo"
echo "  - Recommended steps: 20-50 (default was 4)"
echo "  - Recommended resolution: 1024x1024"
echo "  - CFG Scale: 7.0 (guidance strength)"
echo ""
echo "Note: Generation will be slower but higher quality"
echo "================================================"
