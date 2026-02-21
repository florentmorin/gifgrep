package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/steipete/gifgrep/internal/model"
)

const (
	heypsterBaseURL     = "https://heypster-gif.com"
	heypsterSDKPath     = "/sdk"
	heypsterPageSize    = 10
	heypsterDefaultLang = "en"
)

type heypsterTag struct {
	ID  int    `json:"id"`
	Tag string `json:"tag"`
}

type heypsterGIF struct {
	ID      int           `json:"id"`
	GIFMini *string       `json:"gif_mini"`
	GIF     *string       `json:"gif"`
	H265    string        `json:"h265"`
	Tags    []heypsterTag `json:"tags"`
}

type heypsterGIFPage struct {
	Data []heypsterGIF `json:"data"`
}

func fetchHeypsterV1(query string, opts model.Options) ([]model.Result, error) {
	apiKey := os.Getenv("HEYPSTER_API_KEY")
	if apiKey == "" {
		return nil, errors.New("missing HEYPSTER_API_KEY")
	}

	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}

	tags, err := heypsterSearchTags(apiKey, query)
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return nil, nil
	}

	tag := tags[0]
	title := strings.ReplaceAll(tag.Tag, "-", " ")

	var out []model.Result
	pages := (limit + heypsterPageSize - 1) / heypsterPageSize
	for p := 1; p <= pages; p++ {
		gifs, err := heypsterFetchGIFs(apiKey, tag.ID, p)
		if err != nil {
			return nil, err
		}
		for _, g := range gifs {
			gifURL := heypsterGIFURL(g)
			if gifURL == "" {
				continue
			}
			var tagNames []string
			for _, t := range g.Tags {
				tagNames = append(tagNames, strings.ReplaceAll(t.Tag, "-", " "))
			}
			out = append(out, model.Result{
				ID:         fmt.Sprintf("%d", g.ID),
				Title:      title,
				URL:        gifURL,
				PreviewURL: gifURL,
				Tags:       tagNames,
			})
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(gifs) < heypsterPageSize {
			break
		}
	}
	return out, nil
}

func heypsterSearchTags(apiKey, query string) ([]heypsterTag, error) {
	encoded := url.PathEscape(strings.ToLower(query))
	reqURL := fmt.Sprintf("%s%s/tags/%s/%s", heypsterBaseURL, heypsterSDKPath, encoded, heypsterDefaultLang)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gifgrep")
	req.Header.Set("HEYPSTER-API-KEY", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http %d", resp.StatusCode)
	}

	var tags []heypsterTag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func heypsterFetchGIFs(apiKey string, tagID, page int) ([]heypsterGIF, error) {
	reqURL := fmt.Sprintf("%s%s/gifs-tags/%d?page=%d", heypsterBaseURL, heypsterSDKPath, tagID, page)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gifgrep")
	req.Header.Set("HEYPSTER-API-KEY", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http %d", resp.StatusCode)
	}

	var page1 heypsterGIFPage
	if err := json.NewDecoder(resp.Body).Decode(&page1); err != nil {
		return nil, err
	}
	return page1.Data, nil
}

func heypsterGIFURL(g heypsterGIF) string {
	if g.GIFMini != nil && *g.GIFMini != "" {
		return heypsterBaseURL + "/" + *g.GIFMini
	}
	if g.GIF != nil && *g.GIF != "" {
		return heypsterBaseURL + "/" + *g.GIF
	}
	return ""
}
