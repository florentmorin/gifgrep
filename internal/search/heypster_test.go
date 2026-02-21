package search

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/steipete/gifgrep/internal/model"
	"github.com/steipete/gifgrep/internal/testutil"
)

func TestFetchHeypster(t *testing.T) {
	t.Setenv("HEYPSTER_API_KEY", "test-key")
	gifData := testutil.MakeTestGIF()
	testutil.WithTransport(t, &testutil.FakeTransport{GIFData: gifData}, func() {
		out, err := fetchHeypsterV1("star wars", model.Options{Limit: 1})
		if err != nil {
			t.Fatalf("fetchHeypsterV1 failed: %v", err)
		}
		if len(out) != 1 {
			t.Fatalf("expected 1 result, got %d", len(out))
		}
		if out[0].URL == "" || out[0].PreviewURL == "" {
			t.Fatalf("missing URLs")
		}
		if out[0].ID != "1" {
			t.Fatalf("expected ID '1', got %q", out[0].ID)
		}

		_, err = Search("star wars", model.Options{Limit: 1, Source: "heypster"})
		if err != nil {
			t.Fatalf("Search heypster failed: %v", err)
		}
	})
}

func TestFetchHeypsterMissingKey(t *testing.T) {
	t.Setenv("HEYPSTER_API_KEY", "")
	if _, err := fetchHeypsterV1("star wars", model.Options{Limit: 1}); err == nil {
		t.Fatalf("expected missing key error")
	}
}

type heypsterNoTagsTransport struct{}

func (tr *heypsterNoTagsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.HasPrefix(req.URL.Path, "/sdk/tags/") {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`[]`)),
		}, nil
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("not found")),
	}, nil
}

func TestFetchHeypsterNoTags(t *testing.T) {
	t.Setenv("HEYPSTER_API_KEY", "test-key")
	testutil.WithTransport(t, &heypsterNoTagsTransport{}, func() {
		out, err := fetchHeypsterV1("xyznonexistent", model.Options{Limit: 5})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(out) != 0 {
			t.Fatalf("expected 0 results, got %d", len(out))
		}
	})
}

type heypsterBadJSONTransport struct{}

func (tr *heypsterBadJSONTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader("not-json")),
	}, nil
}

type heypsterStatusTransport struct{}

func (tr *heypsterStatusTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("oops")),
	}, nil
}

func TestFetchHeypsterErrors(t *testing.T) {
	t.Setenv("HEYPSTER_API_KEY", "test-key")
	testutil.WithTransport(t, &heypsterBadJSONTransport{}, func() {
		if _, err := fetchHeypsterV1("star wars", model.Options{Limit: 1}); err == nil {
			t.Fatalf("expected json error")
		}
	})
	testutil.WithTransport(t, &heypsterStatusTransport{}, func() {
		if _, err := fetchHeypsterV1("star wars", model.Options{Limit: 1}); err == nil {
			t.Fatalf("expected status error")
		}
	})
}
