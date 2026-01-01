# Kitty graphics protocol (gifgrep)

gifgrep renders inline previews using the **Kitty graphics protocol**: terminal escape sequences that upload image data and place it into a cell grid.

## What gets sent

The protocol uses an APC-style escape sequence:

```text
ESC _G <params> ; <base64 payload> ESC \
```

gifgrep uploads PNG data (base frame and per-frame PNGs) and then either:

- **Kitty / terminals with native animation:** upload all frames and let the terminal animate.
- **Ghostty:** upload a single frame repeatedly (software playback), since Ghostty currently doesn’t play Kitty animations natively.

## How gifgrep maps to Kitty actions

gifgrep uses these Kitty actions:

- `a=T`: upload a (base) image (PNG)
- `a=f`: append animation frames (PNG) with per-frame delay
- `a=a`: configure animation timing and start playback
- `a=p`: place the image into a cell rectangle
- `a=d`: delete image by id (cleanup)

The payload is base64, chunked (4096 chars) to avoid huge control sequences.

## Terminal support

Works in terminals that implement the Kitty graphics protocol, notably:

- Kitty
- Ghostty (image upload works; gifgrep uses software playback for animation)

Terminals like Apple Terminal and iTerm2 do **not** implement Kitty graphics (iTerm2 has its own inline image protocol; see `docs/iterm.md`).

## Links

- Kitty graphics protocol docs: `https://sw.kovidgoyal.net/kitty/graphics-protocol/`

