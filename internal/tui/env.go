package tui

import (
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

type Env struct {
	In         io.Reader
	Out        io.Writer
	FD         int
	IsTerminal func(int) bool
	MakeRaw    func(int) (*term.State, error)
	Restore    func(int, *term.State) error
	GetSize    func(int) (int, int, error)
	SignalCh   <-chan os.Signal
}

var defaultEnvFn = defaultEnv

func defaultEnv() Env {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	return Env{
		In:         os.Stdin,
		Out:        os.Stdout,
		FD:         int(os.Stdin.Fd()),
		IsTerminal: term.IsTerminal,
		MakeRaw:    term.MakeRaw,
		Restore:    term.Restore,
		GetSize:    term.GetSize,
		SignalCh:   sigs,
	}
}

func SetDefaultEnvForTest(fn func() Env) {
	if fn == nil {
		defaultEnvFn = defaultEnv
		return
	}
	defaultEnvFn = fn
}
