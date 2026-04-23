# Pointegrity logo assets

Red single-dot direction, committed.

## Wordmarks

- `wordmark-red.svg` — primary; full "Pointegrity" with red accent dot on the first `i`
- `wordmark-mono.svg` — ink-only; use when color isn't available (fax, monochrome print, single-color embroidery)

## Standalone Pi mark

- `mark-pi.svg` — square mark for places where the full wordmark doesn't fit (sidebar avatar, app icon, social-share square, etc.)

## Sub-product composite

- `composite-pi-pouch.svg` — template showing the `pi · <product>` pattern; use for places where Pointegrity is the parent context to a product
- `composite-pi-poi.svg` — same pattern for POI

## Favicon set

Generated from `mark-pi.svg` via a headless-Chrome render at precise sizes:

- `favicon.svg` — primary; modern browsers render SVG favicons directly
- `favicon.ico` — legacy fallback; contains 16, 32, 48 px bitmaps packed
- `favicon-16x16.png`, `favicon-32x32.png`, `favicon-96x96.png` — PNG fallbacks
- `apple-touch-icon.png` — 180 × 180, for iOS home screen
- `icon-192.png`, `icon-512.png` — Android / PWA manifest icons

## HTML `<head>` block

Paste into any page that should show the Pointegrity favicon:

```html
<link rel="icon" href="/static/logo/favicon.svg" type="image/svg+xml">
<link rel="icon" href="/static/logo/favicon.ico" sizes="any">
<link rel="icon" href="/static/logo/favicon-32x32.png" sizes="32x32" type="image/png">
<link rel="icon" href="/static/logo/favicon-16x16.png" sizes="16x16" type="image/png">
<link rel="apple-touch-icon" href="/static/logo/apple-touch-icon.png" sizes="180x180">
```

For PWA manifest (if we add one later):

```json
{
  "icons": [
    { "src": "/static/logo/icon-192.png", "sizes": "192x192", "type": "image/png" },
    { "src": "/static/logo/icon-512.png", "sizes": "512x512", "type": "image/png" }
  ]
}
```

## Regenerating the PNGs

When the SVG changes, re-render the rasters from the included
`_render.html` helper. From this directory with `python3 -m http.server`
running on `:4400`:

```bash
for size in 16 32 96 180 192 512; do
  uictl nav "http://localhost:4400/static/logo/_render.html?size=$size"
  uictl wait-ready
  uictl eval "(async () => { await document.fonts.ready; return 'ready'; })()"
  sleep 0.3
  uictl screenshot --viewport ${size}x${size} "static/logo/favicon-${size}x${size}.png"
done
```

(Rename to the final filenames as in this directory.)

Then regenerate `favicon.ico`:

```bash
python3 -c "
from PIL import Image
Image.open('icon-512.png').convert('RGBA').save(
  'favicon.ico', format='ICO', sizes=[(16,16),(32,32),(48,48)]
)"
```

## Note on font dependence

The SVG wordmarks use live `<text>` rendering with Inter via a font
stack. On systems without Inter, fallbacks (system-ui / Segoe UI /
SF Pro) render slightly different metrics — the red accent dot's
`cx` was tuned to Inter. If you need pixel-perfect rendering in an
environment that can't be trusted to have Inter, open the SVG in
Figma or Inkscape and convert the text to outlines before shipping.
The PNG rasters in this directory were generated with real Inter
loaded from Google Fonts, so they're font-locked regardless of
display environment.
