package termcaps

import (
	"os"
	"strings"
)

type InlineProtocol int

const (
	InlineNone InlineProtocol = iota
	InlineKitty
	InlineIterm
)

func (p InlineProtocol) String() string {
	switch p {
	case InlineKitty:
		return "kitty"
	case InlineIterm:
		return "iterm"
	default:
		return "none"
	}
}

func DetectInline(getenv func(string) string) InlineProtocol {
	if getenv == nil {
		getenv = os.Getenv
	}

	switch strings.ToLower(strings.TrimSpace(getenv("GIFGREP_INLINE"))) {
	case "kitty":
		return InlineKitty
	case "iterm", "iterm2":
		return InlineIterm
	case "none", "off", "false", "0":
		return InlineNone
	case "", "auto":
	default:
		return InlineNone
	}

	if strings.TrimSpace(getenv("KITTY_WINDOW_ID")) != "" {
		return InlineKitty
	}

	termProgram := strings.ToLower(getenv("TERM_PROGRAM"))
	if strings.Contains(termProgram, "ghostty") {
		return InlineKitty
	}
	if strings.Contains(termProgram, "iterm") || strings.TrimSpace(getenv("ITERM_SESSION_ID")) != "" {
		return InlineIterm
	}
	if strings.Contains(termProgram, "apple_terminal") {
		return InlineNone
	}

	termEnv := strings.ToLower(getenv("TERM"))
	if strings.Contains(termEnv, "xterm-kitty") || strings.Contains(termEnv, "ghostty") {
		return InlineKitty
	}

	return InlineNone
}
