package model

const AppName = "gifgrep"

var Version = "dev"

type Result struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	URL        string   `json:"url"`
	PreviewURL string   `json:"preview_url"`
	Tags       []string `json:"tags,omitempty"`
	Width      int      `json:"width,omitempty"`
	Height     int      `json:"height,omitempty"`
}

type Options struct {
	TUI        bool
	JSON       bool
	IgnoreCase bool
	Invert     bool
	Regex      bool
	Number     bool
	Limit      int
	Source     string
	Mood       string
	Color      string
}
