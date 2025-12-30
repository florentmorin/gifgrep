package app

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/steipete/gifgrep/internal/model"
	"golang.org/x/term"
)

type helpTableRow struct {
	left  string
	right string
}

func printHelpFor(w io.Writer, opts model.Options, cmd string) {
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	switch cmd {
	case "":
		printRootHelp(w, opts)
	case "search":
		printSearchHelp(w, opts)
	case "tui":
		printTUIHelp(w, opts)
	case "still":
		printStillHelp(w, opts)
	case "sheet":
		printSheetHelp(w, opts)
	default:
		printRootHelp(w, opts)
	}
}

func printRootHelp(w io.Writer, opts model.Options) {
	st := helpStyle{useColor: shouldUseColorForWriter(opts, w)}
	width := termWidth(w)

	st.Println(w, st.Header())
	st.Println(w, "")

	st.Println(w, st.Heading("Usage:"))
	st.Println(w, "  gifgrep <query>                # same as: gifgrep search <query>")
	st.Println(w, "  gifgrep <command> [flags]      # see: gifgrep help <command>")
	st.Println(w, "")

	st.Println(w, st.Heading("Commands:"))
	printHelpTable(w, width, []helpTableRow{
		{st.Flag("search"), "Search and print URLs (default)"},
		{st.Flag("tui"), "Interactive browser with inline preview"},
		{st.Flag("still"), "Extract a single PNG frame"},
		{st.Flag("sheet"), "Generate a contact sheet PNG"},
		{st.Flag("help"), "Show help (help <command>)"},
	})
	st.Println(w, "")

	st.Println(w, st.Heading("Global Flags:"))
	printHelpTable(w, width, []helpTableRow{
		{st.Flag("--color <mode>"), "Color output: auto|always|never (default auto)"},
		{st.Flag("--no-color"), "Alias for --color=never"},
		{st.Flag("-v, --verbose"), "Verbose stderr logs"},
		{st.Flag("-q, --quiet"), "Suppress non-essential stderr output"},
		{st.Flag("-h, --help"), "Show help"},
		{st.Flag("--version"), "Show version"},
	})
	st.Println(w, "")

	st.Println(w, st.Heading("Examples:"))
	st.Println(w, "  gifgrep cats")
	st.Println(w, "  gifgrep search --json cats | jq '.[0].url'")
	st.Println(w, "  gifgrep tui cats")
	st.Println(w, "  gifgrep still cat.gif --at 1.5s -o still.png")
	st.Println(w, "  gifgrep sheet cat.gif --frames 12 --cols 4 -o sheet.png")
}

func printSearchHelp(w io.Writer, opts model.Options) {
	st := helpStyle{useColor: shouldUseColorForWriter(opts, w)}
	width := termWidth(w)

	st.Println(w, st.Header())
	st.Println(w, "")

	st.Println(w, st.Heading("Usage:"))
	st.Println(w, "  gifgrep search [flags] <query>")
	st.Println(w, "  gifgrep [flags] <query>        # default command")
	st.Println(w, "")

	st.Println(w, st.Heading("Output:"))
	st.Println(w, "  Default: one line per GIF: <title>\\t<url>")
	st.Println(w, "  --json: JSON array of results")
	st.Println(w, "")

	st.Println(w, st.Heading("Flags:"))
	printHelpTable(w, width, []helpTableRow{
		{st.Flag("--source <s>"), "Source: auto|tenor|giphy (default auto)"},
		{st.Flag("-m, --max <N>"), "Max results to fetch (default 20)"},
		{st.Flag("-n, --number"), "Prefix lines with 1-based index"},
		{st.Flag("--json"), "JSON output (machine-readable)"},
	})
	st.Println(w, "")

	st.Println(w, st.Heading("Environment:"))
	printHelpTable(w, width, []helpTableRow{
		{st.Flag("TENOR_API_KEY"), "Optional (falls back to Tenor demo key)"},
		{st.Flag("GIPHY_API_KEY"), "Required when using --source giphy"},
	})
	st.Println(w, "")

	st.Println(w, st.Heading("Examples:"))
	st.Println(w, "  gifgrep cats | head -n 5")
	st.Println(w, "  gifgrep search --json cats | jq '.[] | .url'")
	st.Println(w, "  GIPHY_API_KEY=... gifgrep search --source giphy cats")
}

func printTUIHelp(w io.Writer, opts model.Options) {
	st := helpStyle{useColor: shouldUseColorForWriter(opts, w)}
	width := termWidth(w)

	st.Println(w, st.Header())
	st.Println(w, "")

	st.Println(w, st.Heading("Usage:"))
	st.Println(w, "  gifgrep tui [flags] [query]")
	st.Println(w, "")

	st.Println(w, st.Heading("Flags:"))
	printHelpTable(w, width, []helpTableRow{
		{st.Flag("--source <s>"), "Source: auto|tenor|giphy (default auto)"},
		{st.Flag("-m, --max <N>"), "Max results to fetch (default 20)"},
	})
	st.Println(w, "")

	st.Println(w, st.Heading("Keys:"))
	st.Println(w, "  /      edit search")
	st.Println(w, "  ↑↓     select")
	st.Println(w, "  d      download selection")
	st.Println(w, "  q      quit")
	st.Println(w, "")

	st.Println(w, st.Heading("Environment:"))
	printHelpTable(w, width, []helpTableRow{
		{st.Flag("TENOR_API_KEY"), "Optional (falls back to Tenor demo key)"},
		{st.Flag("GIPHY_API_KEY"), "Required when using --source giphy"},
	})
	st.Println(w, "")

	st.Println(w, st.Heading("Examples:"))
	st.Println(w, "  gifgrep tui cats")
}

func printStillHelp(w io.Writer, opts model.Options) {
	st := helpStyle{useColor: shouldUseColorForWriter(opts, w)}
	width := termWidth(w)

	st.Println(w, st.Header())
	st.Println(w, "")

	st.Println(w, st.Heading("Usage:"))
	st.Println(w, "  gifgrep still [flags] <path|url>")
	st.Println(w, "")

	st.Println(w, st.Heading("Flags:"))
	printHelpTable(w, width, []helpTableRow{
		{st.Flag("--at <time>"), "Timestamp (e.g. 1.5s or 1.5)"},
		{st.Flag("-o, --output <path>"), "Output path (default still.png), or '-' for stdout"},
	})
	st.Println(w, "")

	st.Println(w, st.Heading("Examples:"))
	st.Println(w, "  gifgrep still cat.gif --at 0 -o still.png")
	st.Println(w, "  gifgrep still https://example.com/cat.gif --at 1.25s -o - > still.png")
}

func printSheetHelp(w io.Writer, opts model.Options) {
	st := helpStyle{useColor: shouldUseColorForWriter(opts, w)}
	width := termWidth(w)

	st.Println(w, st.Header())
	st.Println(w, "")

	st.Println(w, st.Heading("Usage:"))
	st.Println(w, "  gifgrep sheet [flags] <path|url>")
	st.Println(w, "")

	st.Println(w, st.Heading("Flags:"))
	printHelpTable(w, width, []helpTableRow{
		{st.Flag("--frames <N>"), "Number of frames to sample (default 12)"},
		{st.Flag("--cols <N>"), "Columns (0 = auto)"},
		{st.Flag("--padding <px>"), "Padding between frames (default 2)"},
		{st.Flag("-o, --output <path>"), "Output path (default sheet.png), or '-' for stdout"},
	})
	st.Println(w, "")

	st.Println(w, st.Heading("Examples:"))
	st.Println(w, "  gifgrep sheet cat.gif -o sheet.png")
	st.Println(w, "  gifgrep sheet cat.gif --frames 16 --cols 4 --padding 4 -o sheet.png")
}

func printHelpTable(w io.Writer, termCols int, rows []helpTableRow) {
	if len(rows) == 0 {
		return
	}
	const leftIndent = 2
	const colGap = 2

	leftWidth := 0
	for _, r := range rows {
		if n := visibleLen(r.left); n > leftWidth {
			leftWidth = n
		}
	}

	rightWidth := termCols - leftIndent - leftWidth - colGap
	if rightWidth < 20 {
		rightWidth = 20
	}

	for _, r := range rows {
		left := r.left
		rightLines := wrapText(r.right, rightWidth)
		if len(rightLines) == 0 {
			rightLines = []string{""}
		}

		_, _ = fmt.Fprintf(w, "%s%s%s%s\n",
			strings.Repeat(" ", leftIndent),
			padRightVisible(left, leftWidth),
			strings.Repeat(" ", colGap),
			rightLines[0],
		)
		for _, line := range rightLines[1:] {
			_, _ = fmt.Fprintf(w, "%s%s%s%s\n",
				strings.Repeat(" ", leftIndent),
				strings.Repeat(" ", leftWidth),
				strings.Repeat(" ", colGap),
				line,
			)
		}
	}
}

type helpStyle struct {
	useColor bool
}

func (s helpStyle) Header() string {
	if !s.useColor {
		return fmt.Sprintf("%s %s — %s", model.AppName, model.Version, model.Tagline)
	}
	return "\x1b[1m\x1b[36m" + model.AppName + "\x1b[0m" +
		" " +
		"\x1b[1m" + model.Version + "\x1b[0m" +
		"\x1b[90m — " + model.Tagline + "\x1b[0m"
}

func (s helpStyle) Heading(text string) string {
	if !s.useColor {
		return text
	}
	return "\x1b[1m" + text + "\x1b[0m"
}

func (s helpStyle) Flag(text string) string {
	if !s.useColor {
		return text
	}
	return "\x1b[36m" + text + "\x1b[0m"
}

func (s helpStyle) Println(w io.Writer, text string) {
	_, _ = fmt.Fprintln(w, text)
}

func termWidth(w io.Writer) int {
	if f, ok := w.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		if cols, _, err := term.GetSize(int(f.Fd())); err == nil && cols > 0 {
			return cols
		}
	}
	if raw := strings.TrimSpace(os.Getenv("COLUMNS")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			return v
		}
	}
	return 80
}

func visibleLen(s string) int {
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
		visible++
		i++
	}
	return visible
}

func padRightVisible(s string, width int) string {
	pad := width - visibleLen(s)
	if pad <= 0 {
		return s
	}
	return s + strings.Repeat(" ", pad)
}

func wrapText(s string, width int) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if width <= 0 {
		return []string{s}
	}

	words := strings.Fields(s)
	lines := make([]string, 0, 2)
	var line strings.Builder

	for _, word := range words {
		if line.Len() == 0 {
			line.WriteString(word)
			continue
		}
		if visibleLen(line.String())+1+visibleLen(word) <= width {
			line.WriteByte(' ')
			line.WriteString(word)
			continue
		}
		lines = append(lines, line.String())
		line.Reset()
		line.WriteString(word)
	}
	if line.Len() > 0 {
		lines = append(lines, line.String())
	}
	return lines
}
