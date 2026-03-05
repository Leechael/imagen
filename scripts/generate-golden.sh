#!/bin/bash
# Generate golden reference images using ffmpeg + ImageMagick.
# Run once to create baseline outputs for regression tests.
#
# Prerequisites: ffmpeg, imagemagick (magick)
#
# Usage: bash scripts/generate-golden.sh

set -euo pipefail

TESTDATA="internal/keyer/testdata"
GOLDEN="$TESTDATA/golden"

mkdir -p "$GOLDEN"

for input in "$TESTDATA"/*.png; do
  [[ "$input" == *golden* ]] && continue
  base=$(basename "$input" .png)
  echo "Processing $base..."

  # Detect key color from top-left 4×4 pixels.
  key=$(magick "$input" -crop 4x4+0+0 +repage -format "%c" histogram:info:- \
    | grep -oE '[0-9]+:.*#[0-9A-Fa-f]{6}' \
    | sort -t: -k1 -rn \
    | head -1 \
    | grep -oE '#[0-9A-Fa-f]{6}' \
    | tr -d '#')

  if [ -z "$key" ]; then
    key="00FF00"
  fi

  echo "  key color: 0x${key}"

  # ffmpeg colorkey + despill.
  tmp=$(mktemp /tmp/golden-XXXXXX.png)
  ffmpeg -y -i "$input" -vf "colorkey=0x${key}:0.25:0.08,despill=green" "$tmp" 2>/dev/null

  # magick trim.
  magick "$tmp" -trim +repage "$GOLDEN/${base}.png"
  rm -f "$tmp"

  echo "  → $GOLDEN/${base}.png"
done

echo "Done."
