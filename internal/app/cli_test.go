package app

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/steipete/gifgrep/internal/model"
	"github.com/steipete/gifgrep/internal/testutil"
)

func TestRunSearchOutput(t *testing.T) {
	gifData := testutil.MakeTestGIF()
	testutil.WithTransport(t, &testutil.FakeTransport{GIFData: gifData}, func() {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		t.Cleanup(func() {
			os.Stdout = oldStdout
		})

		err := runSearch(model.Options{Number: true, Limit: 1, Source: "tenor"}, "cats")
		_ = w.Close()
		if err != nil {
			t.Fatalf("runSearch failed: %v", err)
		}
		out, _ := io.ReadAll(r)
		if !strings.Contains(string(out), "1\t") {
			t.Fatalf("expected numbered output")
		}
	})
}

func TestRunSearchJSON(t *testing.T) {
	gifData := testutil.MakeTestGIF()
	testutil.WithTransport(t, &testutil.FakeTransport{GIFData: gifData}, func() {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		t.Cleanup(func() {
			os.Stdout = oldStdout
		})

		err := runSearch(model.Options{JSON: true, Limit: 1, Source: "tenor"}, "cats")
		_ = w.Close()
		if err != nil {
			t.Fatalf("runSearch json failed: %v", err)
		}
		out, _ := io.ReadAll(r)
		if !bytes.Contains(out, []byte(`"preview_url"`)) {
			t.Fatalf("expected json output")
		}
	})
}

func TestHelpOutput(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	t.Cleanup(func() {
		os.Stdout = oldStdout
	})

	code := Run([]string{"--help"})
	_ = w.Close()
	if code != 0 {
		t.Fatalf("expected exit 0")
	}
	out, _ := io.ReadAll(r)
	text := string(out)
	if !strings.Contains(text, "Examples:") {
		t.Fatalf("expected Examples section")
	}
	if !strings.Contains(text, "--no-color") {
		t.Fatalf("expected --no-color in help")
	}
}
