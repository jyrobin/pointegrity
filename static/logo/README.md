# Pointegrity logo assets — FONT-LOCKED

All wordmarks and marks use extracted IBM Plex Sans 700 glyph paths.
No runtime font dependency; identical rendering in every context
(browser, `<img>` tag, LS checkout iframe, PDF embed, email).

## Files

### SVG masters

- `wordmark-red.svg` — full "Pointegrity" + red accent dot (canonical)
- `wordmark-mono.svg` — full "Pointegrity", ink only (single-color contexts)
- `mark-pi.svg` — standalone "pı" + red accent dot (avatar / app icon)
- `favicon.svg` — copy of mark-pi.svg; modern browsers render this in tabs
- `composite-pi-pouch.svg` — `pi · pouch` template for sub-product branding — **still text-based, not font-locked yet** (low priority; relies on Plex being loaded by the host page)
- `composite-pi-poi.svg` — `pi · POI` — same text-based caveat as above

### PNG rasters (generated from mark-pi.svg)

- `favicon.ico` — legacy; 16/32/48 bitmaps packed
- `favicon-16x16.png`, `favicon-32x32.png`, `favicon-96x96.png`
- `apple-touch-icon.png` — 180 × 180, iOS home-screen
- `icon-192.png`, `icon-512.png` — Android / PWA manifest

### Tooling (not user-facing)

- `_extract.py` — Python script that reads `static/fonts/ibm-plex-sans-latin-700-normal.woff2` via fontTools and emits `wordmark-red.svg`, `wordmark-mono.svg`, `mark-pi.svg`
- `_render.html` — tiny page that displays `mark-pi.svg` at a set size; used by headless Chrome to rasterize the favicon PNGs
- `preview.html` — side-by-side display of every committed asset

## Regenerating

If the type or composition needs to change:

```bash
# 1. Edit _extract.py (text, sizes, baseline, viewBox, dot position)
python3 _extract.py

# 2. From the repo root, with python3 -m http.server 4400 running,
#    re-render the PNG favicons
for size in 16 32 96 180 192 512; do
  uictl nav "http://localhost:4400/static/logo/_render.html?size=$size"
  uictl wait-ready
  uictl screenshot --viewport ${size}x${size} \
    "static/logo/<matching-filename>.png"
done

# 3. Regenerate favicon.ico from the 512px master
cd static/logo
python3 -c "
from PIL import Image
Image.open('icon-512.png').convert('RGBA').save(
  'favicon.ico', format='ICO', sizes=[(16,16),(32,32),(48,48)])"
```

## HTML `<head>` block

```html
<link rel="icon" href="/static/logo/favicon.svg" type="image/svg+xml">
<link rel="icon" href="/static/logo/favicon.ico" sizes="any">
<link rel="icon" href="/static/logo/favicon-32x32.png" sizes="32x32" type="image/png">
<link rel="icon" href="/static/logo/favicon-16x16.png" sizes="16x16" type="image/png">
<link rel="apple-touch-icon" href="/static/logo/apple-touch-icon.png" sizes="180x180">
```

## Notes on the typeface

We adopted IBM Plex Sans (bold 700) as the brand typeface — see
`/static/brand.css` and the woff2 files under `/static/fonts/`. The
wordmark and mark were generated from the same 700-weight woff2, so
the printed mark and the site's Plex-rendered chrome are typographically
identical.

If we ever migrate to a different typeface, re-run `_extract.py`
pointing it at the new font's woff2 and the SVG paths update
automatically.
