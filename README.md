# pointegrity

The umbrella site for [Pointegrity](https://www.pointegrity.com) — an
independent software company building privacy-first personal data tools.

## What's here

- `index.html` — landing page (hero, products, philosophy)
- `about/index.html` — `/about` route (company info, why paid-from-day-one, contact)
- `static/motif/` — vendored copy of the motif design tokens + components
- `static/site.css` — Pointegrity-specific overrides on top of motif
- `CNAME` — pins the deploy to `www.pointegrity.com`
- `.nojekyll` — opt out of GitHub Pages' Jekyll processing (we're plain HTML)

Product-specific pages live in product repos. This repo is only for the
company-level surface.

## Local dev

```bash
make dev            # python http.server on :4400
open http://localhost:4400
```

## Deploy (GitHub Pages)

Push to `main`. GitHub Pages is configured to serve from the branch root
(`Settings → Pages → Source: main / root`). Deploy takes ~60 seconds.

### Custom domain

- DNS at the registrar:
  - `www.pointegrity.com` → CNAME → `jyrobin.github.io`
  - `pointegrity.com` (apex) → four GitHub Pages A records:
    `185.199.108.153`, `185.199.109.153`, `185.199.110.153`, `185.199.111.153`
- In GitHub: **Settings → Pages → Custom domain** → `www.pointegrity.com`
  (enforce HTTPS once the cert provisions — usually within an hour)

## Updating vendored motif

If motif ships a token change that we want to pick up:

```bash
make update-motif
git diff static/motif/   # sanity-check
git commit -am "motif: bump vendored CSS"
```

## Adding pages

For a clean URL like `/journal`, create `journal/index.html`. For a
single-page addition like `/terms`, either `terms.html` or `terms/index.html`
works — the `/index.html` variant is cleaner on GH Pages (no extension in URL).

## Design

Colors, spacing, and component styles come from
[motif](https://github.com/jyrobin/motif); overrides live in
`static/site.css`. Keep the overrides small — when a style feels useful for
other sites, push it upstream to motif.
