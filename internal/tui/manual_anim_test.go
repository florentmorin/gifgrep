package tui

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/steipete/gifgrep/gifdecode"
	"github.com/steipete/gifgrep/internal/termcaps"
)

func TestSoftwareAnimationAdvance(t *testing.T) {
	state := &appState{
		useSoftwareAnim: true,
		inline:          termcaps.InlineKitty,
		currentAnim: &gifAnimation{
			ID: 1,
			Frames: []gifdecode.Frame{
				{PNG: []byte{1, 2, 3}, Delay: 10 * time.Millisecond},
				{PNG: []byte{4, 5, 6}, Delay: 10 * time.Millisecond},
			},
		},
		previewNeedsSend: true,
		previewRow:       2,
		previewCol:       2,
	}
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)
	drawPreview(state, out, 10, 5, 2, 2)
	_ = out.Flush()
	if !strings.Contains(buf.String(), "a=T") {
		t.Fatalf("expected initial frame")
	}
	if !state.manualAnim || state.manualNext.IsZero() {
		t.Fatalf("expected manual animation state")
	}

	buf.Reset()
	state.manualNext = time.Now().Add(-time.Millisecond)
	advanceManualAnimation(state, out)
	_ = out.Flush()
	if !strings.Contains(buf.String(), "a=T") {
		t.Fatalf("expected frame update")
	}
}

func TestDrawPreviewSoftwareDirty(t *testing.T) {
	state := &appState{
		useSoftwareAnim: true,
		inline:          termcaps.InlineKitty,
		currentAnim: &gifAnimation{
			ID: 1,
			Frames: []gifdecode.Frame{
				{PNG: []byte{1, 2, 3}, Delay: 10 * time.Millisecond},
			},
		},
		activeImageID:    99,
		previewDirty:     true,
		lastPreview:      struct{ cols, rows int }{cols: 1, rows: 1},
		previewNeedsSend: false,
	}
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)
	drawPreviewSoftware(state, out, 10, 5, 2, 2)
	_ = out.Flush()
	if !strings.Contains(buf.String(), "a=d") {
		t.Fatalf("expected delete for old image")
	}
	if !strings.Contains(buf.String(), "a=T") {
		t.Fatalf("expected redraw")
	}
}

func TestAdvanceManualAnimationGuards(t *testing.T) {
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)
	state := &appState{}
	advanceManualAnimation(state, out)
	if buf.Len() != 0 {
		t.Fatalf("expected no output")
	}

	state.manualAnim = true
	state.currentAnim = &gifAnimation{Frames: []gifdecode.Frame{{PNG: []byte{1}, Delay: 10 * time.Millisecond}}}
	advanceManualAnimation(state, out)

	state.currentAnim = &gifAnimation{Frames: []gifdecode.Frame{{PNG: []byte{1}, Delay: 10 * time.Millisecond}, {PNG: []byte{2}, Delay: 10 * time.Millisecond}}}
	state.lastPreview = struct{ cols, rows int }{cols: 0, rows: 0}
	advanceManualAnimation(state, out)

	state.lastPreview = struct{ cols, rows int }{cols: 10, rows: 5}
	state.manualNext = time.Time{}
	advanceManualAnimation(state, out)

	state.manualNext = time.Now().Add(time.Hour)
	state.previewRow = 1
	state.previewCol = 1
	advanceManualAnimation(state, out)
}
