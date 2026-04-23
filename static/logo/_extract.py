#!/usr/bin/env python3
"""
Extract Plex Sans glyph paths and compose font-locked SVG logos.

This one-shot script reads an IBM Plex Sans woff2 and emits two SVGs
(wordmark + mark) where the text has been converted to <path> elements
— no runtime font dependency. Once generated, the SVGs render pixel-
identical on every platform and in every context (including <img>
tags and LS checkout iframes).

Run from this directory:
    python3 _extract.py

Outputs:
    wordmark-red.svg    (Poıntegrity + red dot)
    mark-pi.svg         (pı + red dot)

Font axes used:
    Plex Sans 700 (Bold) — the heaviest weight Plex ships.
"""
import os
from fontTools.ttLib import TTFont
from fontTools.pens.svgPathPen import SVGPathPen

HERE = os.path.dirname(__file__)
FONTS = os.path.join(HERE, '..', 'fonts')
WOFF2 = os.path.join(FONTS, 'ibm-plex-sans-latin-700-normal.woff2')

# --- helpers ----------------------------------------------------------------

def extract_glyphs(woff2_path, text):
    font = TTFont(woff2_path)
    upem = font['head'].unitsPerEm
    cmap = font.getBestCmap()
    glyphset = font.getGlyphSet()
    hmtx = font['hmtx']
    out = []
    for ch in text:
        gname = cmap.get(ord(ch))
        if not gname:
            raise ValueError(f"char {ch!r} U+{ord(ch):04X} not in font cmap")
        glyph = glyphset[gname]
        pen = SVGPathPen(glyphset)
        glyph.draw(pen)
        advance, lsb = hmtx[gname]
        out.append({'char': ch, 'd': pen.getCommands() or '', 'advance': advance})
    return out, upem

def compose_groups(glyphs, upem, font_size_px, x_start, baseline_y):
    """Return list of <g> strings, one per non-empty glyph, and the end-x."""
    scale = font_size_px / upem
    groups, x = [], x_start
    for g in glyphs:
        if g['d']:
            groups.append(
                f'  <g transform="translate({x:.2f} {baseline_y:.2f}) '
                f'scale({scale:.5f} {-scale:.5f})">'
                f'<path d="{g["d"]}"/></g>'
            )
        x += g['advance'] * scale
    return groups, x

# --- WORDMARK ---------------------------------------------------------------
# Keep the overall composition close to the original (viewBox 560x140,
# 88px type, baseline y=100). Red accent-dot cx is calibrated against
# Plex metrics after extraction.

def build_wordmark(outpath):
    glyphs, upem = extract_glyphs(WOFF2, 'Poıntegrity')
    groups, end_x = compose_groups(glyphs, upem, 88, 20, 100)
    # Locate the dotless ı (index 2 in 'Poıntegrity') to place the red dot.
    scale = 88 / upem
    x = 20
    for i, g in enumerate(glyphs):
        if i == 2:  # ı
            i_left = x
            i_right = x + g['advance'] * scale
            break
        x += g['advance'] * scale
    dot_cx = (i_left + i_right) / 2
    # Plex 'i' dot typically sits ~60-70% of cap-height above baseline.
    # Cap-height in Plex Sans 700 ≈ 698/1000 * 88 ≈ 61.4, so cap-top ≈ 38.6.
    # Dot visually looks best about 5 px above cap-top (i.e. ~y=33).
    # Use r=8 like the user's manual tuning.
    dot_cy = 33
    dot_r = 8
    svg = f'''<?xml version="1.0" encoding="UTF-8"?>
<!--
  Pointegrity wordmark — FONT-LOCKED.
  All glyphs are <path> elements extracted from IBM Plex Sans 700 via
  fontTools. No font-family reference; identical rendering on every
  system, in <img> contexts, in LS checkouts, in email, etc.
  Regenerate: python3 _extract.py (same directory).
-->
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 560 140"
     role="img" aria-label="Pointegrity">
  <g fill="#1e2f14">
{chr(10).join(groups)}
  </g>
  <circle cx="{dot_cx:.2f}" cy="{dot_cy}" r="{dot_r}" fill="#dc2626"/>
</svg>
'''
    with open(outpath, 'w') as f:
        f.write(svg)
    print(f"wrote {outpath}  end_x={end_x:.1f}  dot_cx={dot_cx:.2f}")

# --- MARK (standalone π square) --------------------------------------------

def build_mark_pi(outpath):
    glyphs, upem = extract_glyphs(WOFF2, 'pı')
    # Square viewBox 160x160, font-size 128, baseline around y=120.
    groups, end_x = compose_groups(glyphs, upem, 128, 28, 120)
    # Locate ı
    scale = 128 / upem
    x = 28
    for i, g in enumerate(glyphs):
        if i == 1:  # ı
            i_left = x
            i_right = x + g['advance'] * scale
            break
        x += g['advance'] * scale
    dot_cx = (i_left + i_right) / 2
    dot_cy = 33
    dot_r = 13
    svg = f'''<?xml version="1.0" encoding="UTF-8"?>
<!--
  Pi mark — FONT-LOCKED. Same technique as wordmark.
-->
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 160 160"
     role="img" aria-label="π — Pointegrity">
  <g fill="#1e2f14">
{chr(10).join(groups)}
  </g>
  <circle cx="{dot_cx:.2f}" cy="{dot_cy}" r="{dot_r}" fill="#dc2626"/>
</svg>
'''
    with open(outpath, 'w') as f:
        f.write(svg)
    print(f"wrote {outpath}  end_x={end_x:.1f}  dot_cx={dot_cx:.2f}")

def build_wordmark_mono(outpath):
    """Ink-only wordmark — same glyph paths as wordmark-red, no accent dot."""
    glyphs, upem = extract_glyphs(WOFF2, 'Pointegrity')  # regular 'i' with dot
    groups, end_x = compose_groups(glyphs, upem, 88, 20, 100)
    svg = f'''<?xml version="1.0" encoding="UTF-8"?>
<!--
  Pointegrity wordmark — monochrome, FONT-LOCKED. Use on fax,
  single-color print, embroidery, or any context where the red
  accent can't be reproduced.
-->
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 560 140"
     role="img" aria-label="Pointegrity">
  <g fill="#1e2f14">
{chr(10).join(groups)}
  </g>
</svg>
'''
    with open(outpath, 'w') as f:
        f.write(svg)
    print(f"wrote {outpath}  end_x={end_x:.1f}")

if __name__ == '__main__':
    build_wordmark(os.path.join(HERE, 'wordmark-red.svg'))
    build_wordmark_mono(os.path.join(HERE, 'wordmark-mono.svg'))
    build_mark_pi(os.path.join(HERE, 'mark-pi.svg'))
