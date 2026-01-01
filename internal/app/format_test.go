package app

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/steipete/gifgrep/gifdecode"
	"github.com/steipete/gifgrep/internal/model"
	"github.com/steipete/gifgrep/internal/termcaps"
)

func TestResolveOutputFormatAutoTTY(t *testing.T) {
	prev := isTerminalWriter
	isTerminalWriter = func(_ io.Writer) bool { return true }
	t.Cleanup(func() { isTerminalWriter = prev })

	got := resolveOutputFormat(model.Options{Format: "auto"}, &bytes.Buffer{})
	if got != formatPlain {
		t.Fatalf("expected plain, got %q", got)
	}
}

func TestResolveOutputFormatAutoNonTTY(t *testing.T) {
	prev := isTerminalWriter
	isTerminalWriter = func(_ io.Writer) bool { return false }
	t.Cleanup(func() { isTerminalWriter = prev })

	got := resolveOutputFormat(model.Options{Format: "auto"}, &bytes.Buffer{})
	if got != formatURL {
		t.Fatalf("expected url, got %q", got)
	}
}

func TestRenderPlainNoThumbs(t *testing.T) {
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)

	renderPlain(out, model.Options{Number: true}, false, termcaps.InlineNone, []model.Result{
		{Title: "A dog", URL: "https://example.test/a.gif"},
	})
	_ = out.Flush()

	text := buf.String()
	if !strings.Contains(text, "1. A dog") {
		t.Fatalf("missing title: %q", text)
	}
	if !strings.Contains(text, "  https://example.test/a.gif") {
		t.Fatalf("missing url: %q", text)
	}
}

func TestRenderPlainThumbsNoExtraBlankLine(t *testing.T) {
	prevFetch := fetchThumb
	prevDecode := decodeThumb
	prevSend := sendThumbKitty
	t.Cleanup(func() {
		fetchThumb = prevFetch
		decodeThumb = prevDecode
		sendThumbKitty = prevSend
	})

	fetchThumb = func(_ string) ([]byte, error) { return []byte("gif"), nil }
	decodeThumb = func(_ []byte) (*gifdecode.Frames, error) {
		return &gifdecode.Frames{Frames: []gifdecode.Frame{{PNG: []byte{1}}}}, nil
	}
	sendThumbKitty = func(out *bufio.Writer, id uint32, _ gifdecode.Frame, _, _ int) {
		_, _ = fmt.Fprintf(out, "<IMG%d>", id)
	}

	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)

	renderPlain(out, model.Options{}, false, termcaps.InlineKitty, []model.Result{
		{Title: "A", URL: "https://example.test/a.gif"},
		{Title: "B", URL: "https://example.test/b.gif"},
	})
	_ = out.Flush()

	text := buf.String()
	if !strings.Contains(text, "<IMG1>") || !strings.Contains(text, "<IMG2>") {
		t.Fatalf("expected image markers: %q", text)
	}
	if strings.Contains(text, "\n\n<IMG2>") {
		t.Fatalf("unexpected blank line between thumb blocks: %q", text)
	}
}
