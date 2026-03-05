---
name: nano-banana
description: Generates AI images using the nano-banana CLI (Gemini 3.1 Flash default, Pro available). Handles multi-resolution (512-4K), aspect ratios, reference images for style transfer, green screen workflow for transparent assets, cost tracking, and exact dimension control. Use when asked to "generate an image", "create a sprite", "make an asset", "generate artwork", or any image generation task for UI mockups, game assets, videos, or marketing materials.
---

# nano-banana

AI image generation CLI. Default model: Gemini 3.1 Flash Image Preview.

## Prerequisites

- Go 1.25+
- `GEMINI_API_KEY` environment variable set (or `~/.nano-banana/.env`)
- (Optional) FFmpeg + ImageMagick for transparent mode

Get a Gemini API key at: https://aistudio.google.com/apikey

## /init - First-Time Setup

When the user says "init", "setup nano-banana", or "install nano-banana":

```bash
# 1. Build and install
cd <project-root>
go mod tidy
mkdir -p ~/.local/bin
go build -o ~/.local/bin/nano-banana ./cmd/nano-banana

# 2. Set up API key (environment variable recommended)
export GEMINI_API_KEY=<ask user for their key>

# Or use dotenv file
mkdir -p ~/.nano-banana
echo "GEMINI_API_KEY=<ask user for their key>" > ~/.nano-banana/.env
```

If command not found, ensure PATH contains `~/.local/bin`:
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

## Quick Reference

- Command: `nano-banana "prompt" [options]`
- Default: 1K resolution, Flash model, current directory

## Core Options

| Option | Default | Description |
|--------|---------|-------------|
| `-o, --output` | `nano-gen-{timestamp}` | Output filename (no extension) |
| `-s, --size` | `1K` | Image size: `512`, `1K`, `2K`, or `4K` |
| `-a, --aspect` | model default | Aspect ratio: `1:1`, `16:9`, `9:16`, `4:3`, `3:4`, etc. |
| `-m, --model` | `flash` | Model: `flash`/`nb2`, `pro`/`nb-pro`, or any model ID |
| `-d, --dir` | current directory | Output directory |
| `-r, --ref` | - | Reference image (can use multiple times) |
| `-t, --transparent` | - | Generate on green screen, remove background (FFmpeg) |
| `--api-key` | - | Gemini API key (overrides env/file) |
| `--costs` | - | Show cost summary |
| `--json` | - | JSON output (stdout, script-friendly) |
| `--plain` | - | Plain output (stdout) |
| `--jq EXPR` | - | Filter JSON output with jq expression |

## Models

| Alias | Model | Use When |
|-------|-------|----------|
| `flash`, `nb2` | Gemini 3.1 Flash | Default. Fast, cheap |
| `pro`, `nb-pro` | Gemini 3 Pro | Highest quality needed |

## API Key Resolution Order

1. `--api-key` flag
2. `GEMINI_API_KEY` environment variable
3. `.env` file in current directory
4. `.env` file next to the CLI binary
5. `~/.nano-banana/.env`

## Key Workflows

### Basic Generation

```bash
nano-banana "minimal dashboard UI with dark theme"
nano-banana "cinematic landscape" -s 2K -a 16:9
nano-banana "quick concept sketch" -s 512
```

### Model Selection

```bash
# Default (Flash - fast, cheap)
nano-banana "your prompt"

# Pro (highest quality)
nano-banana "detailed portrait" --model pro -s 2K
```

### Reference Images (Style Transfer / Editing)

```bash
# Edit existing image
nano-banana "change the background to pure white" -r dark-ui.png -o light-ui

# Style transfer - multiple references
nano-banana "combine these two styles" -r style1.png -r style2.png -o combined
```

### Transparent Assets

```bash
nano-banana "robot mascot character" -t -o mascot
nano-banana "pixel art treasure chest" -t -o chest
```

The `-t` flag automatically prompts the AI to generate on a green screen, then uses FFmpeg `colorkey` + `despill` to key out the background and remove green spill from edge pixels.

Requires: `brew install ffmpeg imagemagick`

### Exact Dimensions

1. First `-r` flag: your reference/style image
2. Last `-r` flag: blank image in target dimensions
3. Include dimensions in prompt

```bash
nano-banana "pixel art character in style of first image, 256x256" -r style.png -r blank-256x256.png -o sprite
```

## Cost Tracking

Every generation is logged to `~/.nano-banana/costs.json`. View summary:

```bash
nano-banana --costs
nano-banana --costs --json
```

## Prompt Examples

```bash
# UI mockups
nano-banana "clean SaaS dashboard with analytics charts, white background"

# Widescreen cinematic
nano-banana "cyberpunk cityscape at sunset" -a 16:9 -s 2K

# Pro quality
nano-banana "premium software product hero image" --model pro

# Quick low-res concept
nano-banana "rough sketch of a robot" -s 512

# Game assets with transparency
nano-banana "pixel art treasure chest" -t -o chest

# Portrait aspect ratio
nano-banana "mobile app onboarding screen" -a 9:16
```
