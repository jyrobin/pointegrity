# pointegrity-brand — draft options

Three candidate brand-CSS layers for review. Each is a minimal delta
on top of motif; pick one, rename it to `brand.css`, and it becomes
the shared root override for every Pointegrity product.

## Preview

Local:

```bash
cd ../../..   # repo root
python3 -m http.server 4400
# open http://localhost:4400/static/brand-drafts/preview.html
# flip between ?brand=system | ?brand=inter | ?brand=plex
```

Live (once pushed): `https://www.pointegrity.com/static/brand-drafts/preview.html?brand=<option>`

## The three options at a glance

| | system.css | inter.css | plex.css |
|---|---|---|---|
| **Sans** | system-ui (SF Pro / Segoe UI / Roboto per-OS) | Inter | IBM Plex Sans |
| **Serif** | Georgia / Charter (fallback only) | Georgia (fallback) | **IBM Plex Serif** (first-class) |
| **Mono** | ui-monospace / SF Mono / Menlo (system) | **JetBrains Mono** | **IBM Plex Mono** |
| **Network cost** | 0 KB | ~35 KB (Inter) + ~30 KB (JetBrains Mono) | ~90 KB (3 weights × 3 cuts) |
| **Feels like** | Native / platform-respectful | Indie-SaaS (Linear, Vercel, Stripe) | Precise / technical / editorial-capable |
| **Long-form reading** | adequate | adequate | **strong** (has real serif) |
| **Code rendering** | OK (system mono) | **good** (JetBrains Mono) | **very good** (Plex Mono, coordinated with sans) |
| **Best fit** | Neutral baseline; pouch/bazaar if type isn't the statement | Pure app products (pouch, bazaar) | Mixed content + app (grove / stonecampus) + app (pouch) |

## Adopting one

Pick a file, copy it one level up, and link it from every page:

```bash
cp static/brand-drafts/<pick>.css static/brand.css
```

Then in each page's `<head>`, after motif but before page-specific CSS:

```html
<link rel="stylesheet" href="/static/motif/tokens.css">
<link rel="stylesheet" href="/static/motif/components.css">
<link rel="stylesheet" href="/static/motif/utilities.css">
<link rel="stylesheet" href="/static/brand.css">    <!-- new -->
<link rel="stylesheet" href="/static/site.css">
```

For other Pointegrity products (pouch, bazaar), they link the same
file cross-domain:

```html
<link rel="stylesheet" href="https://www.pointegrity.com/static/brand.css?v=1">
```

## Production note — self-host the fonts

The `inter.css` and `plex.css` drafts use `@import` from Google Fonts
for fast iteration. Before launch, download the woff2 files and
self-host them under `static/fonts/` with `@font-face` declarations
referencing local paths. This:

- Eliminates the Google-Fonts privacy footnote (EU GDPR concerns)
- Removes a runtime dependency
- Keeps the `Server` header owned by you, not Google

## Picking heuristic

Ask yourself:

1. **Will we ever ship content-first surfaces** (blog, journal, tutorials,
   grove/stonecampus learning material)? If yes → **plex**, because
   the matching serif is essentially free and owning a serif later is
   a headache.

2. **Are we committing to a typography-driven brand identity** (like
   Linear's use of Inter or Stripe's custom typeface)? If yes → **inter**.

3. **Otherwise** → **system**. It costs nothing, feels native
   everywhere, and you can migrate to inter or plex later by swapping
   one file.

My read for Pointegrity + the wider family (including grove/stonecampus
under a different umbrella today but stylistically adjacent): **plex**
is the strongest single-file answer. It handles every surface Pointegrity
is likely to ship in the next two years without adding a second
typeface later.

If you prefer not to commit to a licensed typeface today, **system**
is the honest no-regret choice — you can defer the decision without
losing anything.

**inter** is fine but narrower — it buys SaaS aesthetic at the cost of
needing a second font when long-form content arrives.
