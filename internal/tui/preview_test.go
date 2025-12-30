package tui

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/steipete/gifgrep/gifdecode"
)

func TestLoadSelectedImageEdges(t *testing.T) {
	state := &appState{
		results: []gifResult{},
		cache:   map[string]*gifdecode.Frames{},
	}
	loadSelectedImage(state)
	if state.currentAnim != nil {
		t.Fatalf("expected nil animation for empty results")
	}

	state.results = []gifResult{{Title: "no preview"}}
	state.selected = 0
	loadSelectedImage(state)
	if state.currentAnim != nil {
		t.Fatalf("expected nil animation for empty preview url")
	}

	state.cache["https://example.test/preview.gif"] = &gifdecode.Frames{
		Frames: []gifdecode.Frame{{PNG: []byte{1, 2, 3}, Delay: 80 * time.Millisecond}},
		Width:  1,
		Height: 1,
	}
	state.results = []gifResult{{Title: "cached", PreviewURL: "https://example.test/preview.gif"}}
	loadSelectedImage(state)
	if state.currentAnim == nil || !state.previewNeedsSend {
		t.Fatalf("expected cached animation")
	}

	badTransport := &fakeTransport{gifData: []byte("not-a-gif")}
	withTransport(t, badTransport, func() {
		state.cache = map[string]*gifdecode.Frames{}
		state.results = []gifResult{{Title: "bad", PreviewURL: "https://example.test/preview.gif"}}
		state.selected = 0
		loadSelectedImage(state)
		if state.currentAnim != nil {
			t.Fatalf("expected nil animation on decode error")
		}
	})
}

type errTransport struct{}

func (t *errTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("network")
}

func TestFetchGIFError(t *testing.T) {
	withTransport(t, &errTransport{}, func() {
		if _, err := fetchGIF("https://example.test/preview.gif"); err == nil {
			t.Fatalf("expected fetch error")
		}
	})
}
