# Changelog

## Unreleased

### Added
- **Subcommand CLI**: `search`, `tui`, `still`, `sheet` (bare `gifgrep <query...>` remains an alias for `search`).
- **Reveal output**: `--reveal` (and TUI key `f`) reveals the last output file in the file manager.
- **Output controls**: `--quiet` and `--verbose` (with richer structured help).
- **Help output polish**: colorized banner/sections, plus `--no-color` convenience flag.

### Changed
- **Stills workflow**: replaced `--gif/--still/--stills` flags with `still` and `sheet` subcommands (`--at`, `--frames`, `--cols`, `--padding`).

### Removed
- **Local search filters**: dropped `-i`, `-E`, `--mood`, and `-v` (and the related filtering pipeline).

### Fixed
- **TUI hints**: improved key-hint coloring and layout consistency.

### Developer Experience
- **CLI parsing**: migrated to `kong`.

## 0.1.0 - 2025-12-30

### Added
- **CLI search mode**: `gifgrep [flags] <query...>` prints `<title>\t<url>` (or `--json`), with numbered output (`-n`) and `--color` (`auto|always|never`).
- **Providers**: `--source` (`auto`, `tenor`, `giphy`) with `TENOR_API_KEY` and `GIPHY_API_KEY` support (`auto` prefers Giphy when a key is set).
- **Filters**: `-i` ignore-case, `-E` regex over title/tags, `--mood`, `-v` invert vibe, and `-m` max-results.
- **Flags after query**: supports `gifgrep cats --json` (not just `gifgrep --json cats`).
- **Interactive TUI**: `--tui` raw-mode terminal UI with query/browse states, arrow-key navigation, status line, and key hints (plus `d` to download the current selection).
- **Kitty graphics preview**: inline GIF rendering with automatic cleanup; aspect-ratio-aware sizing with `GIFGREP_CELL_ASPECT`.
- **Animation handling**: Kitty animation stream when supported; software re-render fallback (Ghostty auto-detect or `GIFGREP_SOFTWARE_ANIM=1`).
- **Responsive layout**: list + preview side-by-side on wide terminals, stacked preview on narrow terminals.
- **Preview caching**: in-memory cache keyed by preview URL for fast browsing.
- **Stills extraction**: `--gif` (file/URL) with `--still <time>` for a single frame, or `--stills <N>` for a contact sheet; output via `--out` (file or stdout).
- **Giphy attribution**: inline logo display when Giphy is the active source.

### Fixed
- **Frame offsets**: corrected frame-offset math for still/sheet extraction (with regression coverage).

### Developer Experience
- **Formatter + lint**: gofumpt and golangci-lint, with Makefile/justfile targets.
- **Benchmarks**: synthetic and fixture-backed decode benchmarks.
- **Fixtures**: small, licensed GIF corpus with documented sources (`docs/gif-sources.md`).
- **Visual checks**: Ghostty-web screenshot harness (`pnpm snap`).
- **Docs site**: GitHub Pages content in `docs/`.
