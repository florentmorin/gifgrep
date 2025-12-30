package search

import (
	"testing"

	"github.com/steipete/gifgrep/internal/model"
)

func TestFilterResults(t *testing.T) {
	items := []model.Result{
		{Title: "Happy Cat", Tags: []string{"cute", "cat"}},
		{Title: "Angry Dog", Tags: []string{"dog", "angry"}},
		{Title: "Cat Panic", Tags: []string{"cat", "panic"}},
	}

	out, err := FilterResults(items, "cat", model.Options{Regex: true, IgnoreCase: true})
	if err != nil || len(out) != 2 {
		t.Fatalf("regex filter failed: %v len=%d", err, len(out))
	}

	out, err = FilterResults(items, "panic", model.Options{Mood: "panic", IgnoreCase: true})
	if err != nil || len(out) != 1 {
		t.Fatalf("mood filter failed: %v len=%d", err, len(out))
	}

	out, err = FilterResults(items, "panic", model.Options{Mood: "panic", IgnoreCase: true, Invert: true})
	if err != nil || len(out) != 2 {
		t.Fatalf("invert filter failed: %v len=%d", err, len(out))
	}

	_, err = FilterResults(items, "(", model.Options{Regex: true})
	if err == nil {
		t.Fatalf("expected regex error")
	}
}
