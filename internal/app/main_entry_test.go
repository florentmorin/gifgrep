package app

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/steipete/gifgrep/internal/tui"
	"golang.org/x/term"
)

func TestRunArgs(t *testing.T) {
	t.Run("version", func(t *testing.T) {
		if code := Run([]string{"--version"}); code != 0 {
			t.Fatalf("expected exit 0")
		}
	})

	t.Run("help", func(t *testing.T) {
		if code := Run([]string{"--help"}); code != 0 {
			t.Fatalf("expected exit 0")
		}
	})

	t.Run("empty", func(t *testing.T) {
		if code := Run(nil); code != 1 {
			t.Fatalf("expected exit 1")
		}
	})

	t.Run("bad args", func(t *testing.T) {
		if code := Run([]string{"--nope"}); code != 1 {
			t.Fatalf("expected exit 1")
		}
	})

	t.Run("bad source", func(t *testing.T) {
		if code := Run([]string{"--source", "nope", "cats"}); code != 1 {
			t.Fatalf("expected exit 1")
		}
	})

	t.Run("tui", func(t *testing.T) {
		t.Cleanup(func() { tui.SetDefaultEnvForTest(nil) })
		tui.SetDefaultEnvForTest(func() tui.Env {
			return tui.Env{
				In:         bytes.NewReader([]byte("q")),
				Out:        io.Discard,
				FD:         1,
				IsTerminal: func(int) bool { return true },
				MakeRaw:    func(int) (*term.State, error) { return &term.State{}, nil },
				Restore:    func(int, *term.State) error { return nil },
				GetSize:    func(int) (int, int, error) { return 80, 24, nil },
				SignalCh:   make(chan os.Signal),
			}
		})
		if code := Run([]string{"--tui"}); code != 0 {
			t.Fatalf("expected exit 0")
		}
	})
}

func TestRunTUIExitCodes(t *testing.T) {
	t.Cleanup(func() { tui.SetDefaultEnvForTest(nil) })

	t.Run("not terminal", func(t *testing.T) {
		tui.SetDefaultEnvForTest(func() tui.Env {
			return tui.Env{
				In:         bytes.NewReader(nil),
				Out:        io.Discard,
				FD:         1,
				IsTerminal: func(int) bool { return false },
			}
		})
		if code := Run([]string{"--tui"}); code != 1 {
			t.Fatalf("expected exit 1")
		}
	})

	t.Run("makeRaw fails", func(t *testing.T) {
		tui.SetDefaultEnvForTest(func() tui.Env {
			return tui.Env{
				In:         bytes.NewReader(nil),
				Out:        io.Discard,
				FD:         1,
				IsTerminal: func(int) bool { return true },
				MakeRaw: func(int) (*term.State, error) {
					return nil, errors.New("boom")
				},
			}
		})
		if code := Run([]string{"--tui"}); code != 1 {
			t.Fatalf("expected exit 1")
		}
	})
}
