package iterm

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestSendInlineFile(t *testing.T) {
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)
	SendInlineFile(out, File{
		Name:        "thumb.png",
		Data:        []byte{1, 2, 3},
		WidthCells:  10,
		HeightCells: 5,
	})
	_ = out.Flush()
	s := buf.String()
	if !strings.HasPrefix(s, "\x1b]1337;File=") {
		t.Fatalf("expected osc 1337 prefix")
	}
	if !strings.Contains(s, "inline=1") {
		t.Fatalf("expected inline=1")
	}
	if !strings.Contains(s, "width=10") || !strings.Contains(s, "height=5") {
		t.Fatalf("expected width/height in cells")
	}
	if !strings.HasSuffix(s, "\x1b\\") {
		t.Fatalf("expected ST terminator")
	}
}
