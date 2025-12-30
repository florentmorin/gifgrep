package tui

import (
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func truncateRunes(s string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= width {
		return s
	}
	return string(runes[:width])
}

func truncateANSI(s string, width int) string {
	if width <= 0 {
		return ""
	}
	var out strings.Builder
	out.Grow(len(s))
	visible := 0
	hadANSI := false
	for i := 0; i < len(s); {
		if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && s[j] != 'm' {
				j++
			}
			if j < len(s) {
				j++
			}
			out.WriteString(s[i:j])
			hadANSI = true
			i = j
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			size = 1
		}
		if visible+1 > width {
			break
		}
		out.WriteRune(r)
		visible++
		i += size
	}
	if hadANSI && !strings.HasSuffix(out.String(), "\x1b[0m") {
		out.WriteString("\x1b[0m")
	}
	return out.String()
}

func visibleRuneLen(s string) int {
	visible := 0
	for i := 0; i < len(s); {
		if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && s[j] != 'm' {
				j++
			}
			if j < len(s) {
				j++
			}
			i = j
			continue
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		if size == 0 {
			break
		}
		visible++
		i += size
	}
	return visible
}

func runeLen(s string) int {
	return len([]rune(s))
}

func clampDelay(delay time.Duration) time.Duration {
	if delay < 10*time.Millisecond {
		return 10 * time.Millisecond
	}
	if delay > time.Second {
		return time.Second
	}
	return delay
}

func delayMS(delay time.Duration) int {
	return int(clampDelay(delay).Milliseconds())
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func cellAspectRatio() float64 {
	if raw := strings.TrimSpace(os.Getenv("GIFGREP_CELL_ASPECT")); raw != "" {
		if v, err := strconv.ParseFloat(raw, 64); err == nil && v > 0.1 && v < 2 {
			return v
		}
	}
	return 0.5
}

func useSoftwareAnimation() bool {
	if raw := strings.TrimSpace(os.Getenv("GIFGREP_SOFTWARE_ANIM")); raw != "" {
		raw = strings.ToLower(raw)
		return raw == "1" || raw == "true" || raw == "yes"
	}
	termProgram := strings.ToLower(os.Getenv("TERM_PROGRAM"))
	term := strings.ToLower(os.Getenv("TERM"))
	if strings.Contains(termProgram, "ghostty") || strings.Contains(term, "ghostty") {
		return true
	}
	return false
}

func styleIf(enabled bool, text string, codes ...string) string {
	if !enabled || len(codes) == 0 {
		return text
	}
	var b strings.Builder
	for _, code := range codes {
		b.WriteString(code)
	}
	b.WriteString(text)
	b.WriteString("\x1b[0m")
	return b.String()
}
