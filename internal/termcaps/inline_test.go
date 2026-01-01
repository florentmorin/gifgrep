package termcaps

import "testing"

func TestDetectInlineOverride(t *testing.T) {
	getenv := func(k string) string {
		switch k {
		case "GIFGREP_INLINE":
			return "iterm"
		default:
			return ""
		}
	}
	if got := DetectInline(getenv); got != InlineIterm {
		t.Fatalf("expected iterm, got %v", got)
	}
}

func TestDetectInlineKittyEnv(t *testing.T) {
	getenv := func(k string) string {
		switch k {
		case "KITTY_WINDOW_ID":
			return "123"
		case "TERM_PROGRAM":
			return "iTerm.app"
		default:
			return ""
		}
	}
	if got := DetectInline(getenv); got != InlineKitty {
		t.Fatalf("expected kitty, got %v", got)
	}
}

func TestDetectInlineItermEnv(t *testing.T) {
	getenv := func(k string) string {
		switch k {
		case "TERM_PROGRAM":
			return "iTerm.app"
		default:
			return ""
		}
	}
	if got := DetectInline(getenv); got != InlineIterm {
		t.Fatalf("expected iterm, got %v", got)
	}
}
