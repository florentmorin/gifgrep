package app

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/steipete/gifgrep/internal/model"
	"github.com/steipete/gifgrep/internal/search"
	"golang.org/x/term"
)

var (
	errHelp    = errors.New("help")
	errVersion = errors.New("version")
)

type usageError struct {
	cmd string
	msg string
}

func (e usageError) Error() string {
	if e.msg == "" {
		return "usage error"
	}
	return e.msg
}

func parseArgs(args []string) (string, model.Options, string, error) {
	opts := model.Options{
		Color:         "auto",
		Limit:         20,
		Source:        "auto",
		StillsPadding: 2,
	}

	rest, showHelp, showVersion, err := stripGlobalFlags(&opts, args)
	if err != nil {
		return "", opts, "", usageError{cmd: "", msg: err.Error()}
	}

	if showVersion {
		_, _ = fmt.Fprintf(os.Stdout, "%s %s\n", model.AppName, model.Version)
		return "", opts, "", errVersion
	}

	if len(rest) == 0 {
		if showHelp {
			printHelpFor(os.Stdout, opts, "")
			return "", opts, "", errHelp
		}
		return "", opts, "", usageError{cmd: "", msg: "missing query"}
	}

	cmd, cmdArgs := detectCommand(rest)
	if cmd == "help" {
		target := ""
		if len(cmdArgs) > 0 {
			target, _ = normalizeCommand(cmdArgs[0])
		}
		printHelpFor(os.Stdout, opts, target)
		return "", opts, "", errHelp
	}

	if showHelp {
		printHelpFor(os.Stdout, opts, cmd)
		return "", opts, "", errHelp
	}

	switch cmd {
	case "search":
		query, err := parseSearchArgs(&opts, cmdArgs)
		if err != nil {
			return "", opts, "", err
		}
		return "search", opts, query, nil
	case "tui":
		query, err := parseTUIArgs(&opts, cmdArgs)
		if err != nil {
			return "", opts, "", err
		}
		return "tui", opts, query, nil
	case "still":
		if err := parseStillArgs(&opts, cmdArgs); err != nil {
			return "", opts, "", err
		}
		return "still", opts, "", nil
	case "sheet":
		if err := parseSheetArgs(&opts, cmdArgs); err != nil {
			return "", opts, "", err
		}
		return "sheet", opts, "", nil
	default:
		return "", opts, "", usageError{cmd: "", msg: fmt.Sprintf("unknown command: %s", cmd)}
	}
}

func stripGlobalFlags(opts *model.Options, args []string) ([]string, bool, bool, error) {
	if opts == nil {
		return args, false, false, errors.New("missing opts")
	}

	var showHelp bool
	var showVersion bool
	rest := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			rest = append(rest, args[i+1:]...)
			break
		}

		switch {
		case arg == "--help" || arg == "-h":
			showHelp = true
		case arg == "--version":
			showVersion = true
		case arg == "--no-color":
			opts.Color = "never"
		case strings.HasPrefix(arg, "--color="):
			opts.Color = strings.TrimPrefix(arg, "--color=")
		case arg == "--color":
			if i+1 >= len(args) {
				return nil, false, false, errors.New("missing value for --color")
			}
			i++
			opts.Color = args[i]
		case arg == "-v" || arg == "--verbose":
			opts.Verbose++
		case arg == "-q" || arg == "--quiet":
			opts.Quiet = true
		default:
			rest = append(rest, arg)
		}
	}

	opts.Color = strings.ToLower(strings.TrimSpace(opts.Color))
	switch opts.Color {
	case "", "auto":
		opts.Color = "auto"
	case "always", "never":
	default:
		return nil, false, false, fmt.Errorf("invalid --color: %q (expected auto|always|never)", opts.Color)
	}

	return rest, showHelp, showVersion, nil
}

func detectCommand(args []string) (string, []string) {
	if len(args) == 0 {
		return "help", nil
	}
	cmd, ok := normalizeCommand(args[0])
	if ok {
		return cmd, args[1:]
	}
	return "search", args
}

func normalizeCommand(raw string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "search":
		return "search", true
	case "tui", "browse":
		return "tui", true
	case "still":
		return "still", true
	case "sheet", "contact-sheet", "contactsheet", "stills":
		return "sheet", true
	case "help":
		return "help", true
	default:
		return "", false
	}
}

func parseSearchArgs(opts *model.Options, args []string) (string, error) {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&opts.JSON, "json", false, "json output")
	fs.BoolVar(&opts.Number, "n", false, "number results")
	fs.BoolVar(&opts.Number, "number", false, "number results")
	fs.IntVar(&opts.Limit, "m", opts.Limit, "max results")
	fs.IntVar(&opts.Limit, "max", opts.Limit, "max results")
	fs.IntVar(&opts.Limit, "limit", opts.Limit, "max results")
	fs.StringVar(&opts.Source, "source", opts.Source, "source: auto|tenor|giphy")

	if err := fs.Parse(args); err != nil {
		return "", usageError{cmd: "search", msg: "bad args"}
	}
	if opts.Limit < 1 {
		return "", usageError{cmd: "search", msg: "bad args: --max must be >= 1"}
	}
	if !isValidSource(opts.Source) {
		return "", usageError{cmd: "search", msg: "bad args: --source must be auto|tenor|giphy"}
	}

	query := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if query == "" {
		return "", usageError{cmd: "search", msg: "missing query"}
	}
	return query, nil
}

func parseTUIArgs(opts *model.Options, args []string) (string, error) {
	fs := flag.NewFlagSet("tui", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.IntVar(&opts.Limit, "m", opts.Limit, "max results")
	fs.IntVar(&opts.Limit, "max", opts.Limit, "max results")
	fs.IntVar(&opts.Limit, "limit", opts.Limit, "max results")
	fs.StringVar(&opts.Source, "source", opts.Source, "source: auto|tenor|giphy")

	if err := fs.Parse(args); err != nil {
		return "", usageError{cmd: "tui", msg: "bad args"}
	}
	if opts.Limit < 1 {
		return "", usageError{cmd: "tui", msg: "bad args: --max must be >= 1"}
	}
	if !isValidSource(opts.Source) {
		return "", usageError{cmd: "tui", msg: "bad args: --source must be auto|tenor|giphy"}
	}

	query := strings.TrimSpace(strings.Join(fs.Args(), " "))
	return query, nil
}

func parseStillArgs(opts *model.Options, args []string) error {
	var inputFlag string
	var atRaw string
	var outPath string

	inputPos := ""
	if len(args) > 0 && args[0] != "--" && !strings.HasPrefix(args[0], "-") {
		inputPos = args[0]
		args = args[1:]
	}

	fs := flag.NewFlagSet("still", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&inputFlag, "gif", "", "gif input path or URL")
	fs.StringVar(&atRaw, "at", "", "timestamp (e.g. 1.5s)")
	fs.StringVar(&outPath, "o", "", "output path or '-' for stdout")
	fs.StringVar(&outPath, "output", "", "output path or '-' for stdout")

	if err := fs.Parse(args); err != nil {
		return usageError{cmd: "still", msg: "bad args"}
	}

	input := strings.TrimSpace(inputFlag)
	if input == "" {
		input = strings.TrimSpace(inputPos)
	}
	pos := fs.Args()
	if input == "" && len(pos) > 0 {
		input = pos[0]
		pos = pos[1:]
	}
	if input == "" {
		return usageError{cmd: "still", msg: "missing GIF input"}
	}
	if len(pos) > 0 {
		return usageError{cmd: "still", msg: "unexpected args"}
	}
	if strings.TrimSpace(atRaw) == "" {
		return usageError{cmd: "still", msg: "missing --at"}
	}
	at, err := parseDurationValue(atRaw)
	if err != nil {
		return usageError{cmd: "still", msg: "bad args: invalid --at"}
	}

	opts.GifInput = input
	opts.StillSet = true
	opts.StillAt = at
	opts.OutPath = outPath
	if strings.TrimSpace(opts.OutPath) == "" {
		opts.OutPath = "still.png"
	}
	return nil
}

func parseSheetArgs(opts *model.Options, args []string) error {
	var inputFlag string
	var outPath string
	frames := 12

	inputPos := ""
	if len(args) > 0 && args[0] != "--" && !strings.HasPrefix(args[0], "-") {
		inputPos = args[0]
		args = args[1:]
	}

	fs := flag.NewFlagSet("sheet", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&inputFlag, "gif", "", "gif input path or URL")
	fs.IntVar(&frames, "frames", frames, "frame count")
	fs.IntVar(&opts.StillsCols, "cols", opts.StillsCols, "columns")
	fs.IntVar(&opts.StillsPadding, "padding", opts.StillsPadding, "padding (px)")
	fs.StringVar(&outPath, "o", "", "output path or '-' for stdout")
	fs.StringVar(&outPath, "output", "", "output path or '-' for stdout")

	if err := fs.Parse(args); err != nil {
		return usageError{cmd: "sheet", msg: "bad args"}
	}

	input := strings.TrimSpace(inputFlag)
	if input == "" {
		input = strings.TrimSpace(inputPos)
	}
	pos := fs.Args()
	if input == "" && len(pos) > 0 {
		input = pos[0]
		pos = pos[1:]
	}
	if input == "" {
		return usageError{cmd: "sheet", msg: "missing GIF input"}
	}
	if len(pos) > 0 {
		return usageError{cmd: "sheet", msg: "unexpected args"}
	}
	if frames < 1 {
		return usageError{cmd: "sheet", msg: "bad args: --frames must be >= 1"}
	}
	if opts.StillsCols < 0 {
		return usageError{cmd: "sheet", msg: "bad args: --cols must be >= 0"}
	}
	if opts.StillsPadding < 0 {
		return usageError{cmd: "sheet", msg: "bad args: --padding must be >= 0"}
	}

	opts.GifInput = input
	opts.StillSet = false
	opts.StillsCount = frames
	opts.OutPath = outPath
	if strings.TrimSpace(opts.OutPath) == "" {
		opts.OutPath = "sheet.png"
	}
	return nil
}

func isValidSource(source string) bool {
	source = strings.ToLower(strings.TrimSpace(source))
	switch source {
	case "", "auto", "tenor", "giphy":
		return true
	default:
		return false
	}
}

func runSearch(opts model.Options, query string) error {
	if strings.TrimSpace(query) == "" {
		return errors.New("missing query")
	}
	if opts.Verbose > 0 && !opts.Quiet {
		_, _ = fmt.Fprintf(os.Stderr, "source=%s max=%d\n", search.ResolveSource(opts.Source), opts.Limit)
	}

	results, err := search.Search(query, opts)
	if err != nil {
		return err
	}

	if opts.JSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	useColor := shouldUseColorForWriter(opts, os.Stdout)
	for i, res := range results {
		prefix := ""
		if opts.Number {
			prefix = fmt.Sprintf("%d\t", i+1)
		}
		label := strings.Join(strings.Fields(res.Title), " ")
		if label == "" {
			label = strings.Join(strings.Fields(res.ID), " ")
		}
		if label == "" {
			label = "untitled"
		}
		url := res.URL
		if useColor {
			label = "\x1b[1m" + label + "\x1b[0m"
			url = "\x1b[36m" + url + "\x1b[0m"
		}
		_, _ = fmt.Fprintf(os.Stdout, "%s%s\t%s\n", prefix, label, url)
	}
	return nil
}

func shouldUseColorForWriter(opts model.Options, w io.Writer) bool {
	if opts.Color == "never" {
		return false
	}
	if opts.Color == "always" {
		return true
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	termEnv := strings.ToLower(strings.TrimSpace(os.Getenv("TERM")))
	if termEnv == "dumb" || termEnv == "" {
		return false
	}
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}

func parseDurationValue(raw string) (time.Duration, error) {
	if raw == "" {
		return 0, errors.New("empty duration")
	}
	if d, err := time.ParseDuration(raw); err == nil {
		return d, nil
	}
	if secs, err := strconv.ParseFloat(raw, 64); err == nil {
		if secs < 0 {
			return 0, errors.New("negative duration")
		}
		return time.Duration(secs * float64(time.Second)), nil
	}
	return 0, errors.New("invalid duration")
}
