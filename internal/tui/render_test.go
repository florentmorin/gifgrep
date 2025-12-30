package tui

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/steipete/gifgrep/gifdecode"
	"github.com/steipete/gifgrep/internal/model"
)

func TestRenderDeletesOldImage(t *testing.T) {
	state := &appState{
		mode:          modeBrowse,
		results:       []model.Result{{Title: "A"}},
		activeImageID: 7,
		currentAnim:   nil,
		opts:          model.Options{Source: "tenor"},
	}
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)
	render(state, out, 10, 60)
	_ = out.Flush()
	if !strings.Contains(buf.String(), "a=d") {
		t.Fatalf("expected delete kitty image")
	}
}

func TestRenderWithPreviewRight(t *testing.T) {
	state := &appState{
		mode:    modeBrowse,
		results: []model.Result{{Title: "A"}},
		currentAnim: &gifAnimation{
			ID:     1,
			Frames: []gifdecode.Frame{{PNG: []byte{1, 2, 3}, Delay: 80 * time.Millisecond}},
			Width:  200,
			Height: 100,
		},
		previewNeedsSend: true,
		opts:             model.Options{Source: "tenor"},
	}
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)
	render(state, out, 20, 90)
	_ = out.Flush()
	if !strings.Contains(buf.String(), "a=T") {
		t.Fatalf("expected kitty image data")
	}
}

func TestRenderWithPreviewBottom(t *testing.T) {
	state := &appState{
		mode:    modeBrowse,
		results: []model.Result{{Title: "A"}},
		currentAnim: &gifAnimation{
			ID:     2,
			Frames: []gifdecode.Frame{{PNG: []byte{1, 2, 3}, Delay: 80 * time.Millisecond}},
			Width:  200,
			Height: 100,
		},
		previewNeedsSend: true,
		opts:             model.Options{Source: "tenor"},
	}
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)
	render(state, out, 24, 60)
	_ = out.Flush()
	if !strings.Contains(buf.String(), "Preview") {
		t.Fatalf("expected preview label")
	}
	if !strings.Contains(buf.String(), "a=T") {
		t.Fatalf("expected kitty image data")
	}
}

func TestDrawPreviewPlacement(t *testing.T) {
	state := &appState{
		currentAnim: &gifAnimation{
			ID:     3,
			Frames: []gifdecode.Frame{{PNG: []byte{1, 2, 3}, Delay: 80 * time.Millisecond}},
		},
		previewNeedsSend: false,
		previewDirty:     false,
		activeImageID:    3,
		opts:             model.Options{Source: "tenor"},
		lastPreview: struct {
			cols int
			rows int
		}{cols: 1, rows: 1},
	}
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)
	drawPreview(state, out, 10, 6, 2, 2)
	_ = out.Flush()
	if !strings.Contains(buf.String(), "a=p") {
		t.Fatalf("expected kitty placement")
	}
}

func TestWriteLineAtClears(t *testing.T) {
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)
	writeLineAt(out, 1, 1, "hi", 0)
	_ = out.Flush()
	if !strings.Contains(buf.String(), "\x1b[K") {
		t.Fatalf("expected clear line")
	}
}
