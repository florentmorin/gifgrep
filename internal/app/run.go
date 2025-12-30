package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/steipete/gifgrep/internal/tui"
)

func Run(args []string) int {
	cmd, opts, query, err := parseArgs(args)
	if err != nil {
		if errors.Is(err, errHelp) || errors.Is(err, errVersion) {
			return 0
		}
		var usage usageError
		if errors.As(err, &usage) {
			if usage.msg != "" {
				_, _ = fmt.Fprintln(os.Stderr, usage.msg)
				_, _ = fmt.Fprintln(os.Stderr, "")
			}
			printHelpFor(os.Stderr, opts, usage.cmd)
			return 2
		}
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	switch cmd {
	case "search":
		if err := runSearch(opts, query); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		return 0
	case "tui":
		if err := tui.Run(opts, query); err != nil {
			if errors.Is(err, tui.ErrNotTerminal) {
				_, _ = fmt.Fprintln(os.Stderr, "stdin is not a tty")
			} else {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}
			return 1
		}
		return 0
	case "still":
		if err := runExtract(opts); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		return 0
	case "sheet":
		if err := runExtract(opts); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		return 0
	default:
		_, _ = fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printHelpFor(os.Stderr, opts, "")
		return 2
	}
}
