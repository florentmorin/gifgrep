#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
tmp_dir="$(mktemp -d)"
bin="$tmp_dir/termcaps-check"

cleanup() { rm -rf "$tmp_dir"; }
trap cleanup EXIT

cd "$repo_root"
go build -o "$bin" ./cmd/termcaps-check

wait_for_file() {
  local f="$1"
  local tries=0
  until [[ -s "$f" ]]; do
    tries=$((tries + 1))
    if [[ $tries -gt 80 ]]; then
      return 1
    fi
    sleep 0.1
  done
}

app_bundle_installed() {
  local bundle_id="$1"
  shift
  local paths=("$@")

  for p in "${paths[@]}"; do
    if [[ -d "$p" ]]; then
      return 0
    fi
  done

  if command -v mdfind >/dev/null 2>&1; then
    local found
    found="$(mdfind "kMDItemCFBundleIdentifier == '$bundle_id'" | head -n 1 || true)"
    [[ -n "$found" ]]
    return
  fi

  return 1
}

read_status_or_fail() {
  local label="$1"
  local status_file="$2"
  local out_file="$3"

  if ! wait_for_file "$status_file"; then
    echo "FAIL $label (no status file)"
    [[ -f "$out_file" ]] && cat "$out_file" || true
    exit 1
  fi

  local status
  status="$(tr -d '\r\n ' <"$status_file" || true)"
  if [[ "$status" != "0" ]]; then
    echo "FAIL $label (exit=$status)"
    [[ -f "$out_file" ]] && cat "$out_file" || true
    exit 1
  fi
}

run_terminal_applescript() {
  local app="$1"
  local expect="$2"
  local out="$3"
  local status="$4"

  /usr/bin/osascript <<'APPLESCRIPT' "$app" "$bin" "$expect" "$out" "$status"
on run argv
  set appName to item 1 of argv
  set toolPath to item 2 of argv
  set expectProto to item 3 of argv
  set outPath to item 4 of argv
  set statusPath to item 5 of argv

  set cmd to "/bin/zsh -lc " & quoted form of ("GIFGREP_INLINE=auto " & toolPath & " --expect " & expectProto & " --json > " & outPath & " 2>&1; echo $? > " & statusPath & "; exit")

  if appName is "Terminal" then
    tell application "Terminal"
      activate
      set t to do script cmd
    end tell
  else if appName is "iTerm" then
    tell application "iTerm"
      activate
      set w to (create window with default profile)
      tell current session of w
        write text cmd
      end tell
    end tell
  else
    error "unsupported app: " & appName
  end if
end run
APPLESCRIPT
}

pass() { printf "PASS %s\n" "$1"; }
fail() { printf "FAIL %s\n" "$1"; exit 1; }
skip() { printf "SKIP %s\n" "$1"; }

echo "termcaps e2e (macOS)"

# Current terminal (sanity)
current_out="$tmp_dir/current.json"
if "$bin" --json >"$current_out"; then
  pass "current-terminal"
else
  fail "current-terminal"
fi

# Apple Terminal (expect none)
if app_bundle_installed "com.apple.Terminal" \
  "/System/Applications/Utilities/Terminal.app" \
  "/Applications/Utilities/Terminal.app"; then
  term_out="$tmp_dir/terminal.json"
  term_status="$tmp_dir/terminal.status"
  if run_terminal_applescript "Terminal" "none" "$term_out" "$term_status"; then
    read_status_or_fail "Apple Terminal" "$term_status" "$term_out"
    pass "Apple Terminal"
  else
    skip "Apple Terminal (osascript failed)"
  fi
else
  skip "Apple Terminal (not installed)"
fi

# iTerm2 (expect iterm) — skip if not installed
if app_bundle_installed "com.googlecode.iterm2" \
  "/Applications/iTerm.app" \
  "/Applications/iTerm2.app" \
  "$HOME/Applications/iTerm.app" \
  "$HOME/Applications/iTerm2.app"; then
  iterm_out="$tmp_dir/iterm.json"
  iterm_status="$tmp_dir/iterm.status"
  if run_terminal_applescript "iTerm" "iterm" "$iterm_out" "$iterm_status"; then
    read_status_or_fail "iTerm2" "$iterm_status" "$iterm_out"
    pass "iTerm2"
  else
    fail "iTerm2 (osascript failed)"
  fi
else
  skip "iTerm2 (not installed)"
fi

run_cli_terminal() {
  local label="$1"
  local exe="$2"
  local expect="$3"
  local out="$4"
  local status="$5"
  shift 5

  local cmd
  cmd="GIFGREP_INLINE=auto \"$bin\" --expect \"$expect\" --json >\"$out\" 2>&1; echo \$? >\"$status\""

  "$exe" "$@" /bin/zsh -lc "$cmd" >/dev/null 2>&1 &
  read_status_or_fail "$label" "$status" "$out"
  pass "$label"
}

# Kitty (expect kitty) — skip if not installed
if command -v kitty >/dev/null 2>&1; then
  kitty_out="$tmp_dir/kitty.json"
  kitty_status="$tmp_dir/kitty.status"
  run_cli_terminal "Kitty" "$(command -v kitty)" "kitty" "$kitty_out" "$kitty_status" -e
else
  skip "Kitty (not installed)"
fi

# Ghostty (expect kitty) — skip if CLI not available / unknown flags
if command -v ghostty >/dev/null 2>&1; then
  # Try common flags; if unsupported, skip.
  if ghostty --help 2>&1 | grep -q -- " -e"; then
    ghostty_out="$tmp_dir/ghostty.json"
    ghostty_status="$tmp_dir/ghostty.status"
    run_cli_terminal "Ghostty" "$(command -v ghostty)" "kitty" "$ghostty_out" "$ghostty_status" -e
  else
    skip "Ghostty (no -e flag)"
  fi
else
  skip "Ghostty (not installed)"
fi

echo "Outputs:"
ls -1 "$tmp_dir"/*.json 2>/dev/null || true
