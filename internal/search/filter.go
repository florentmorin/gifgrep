package search

import (
	"regexp"
	"strings"

	"github.com/steipete/gifgrep/internal/model"
)

func FilterResults(results []model.Result, query string, opts model.Options) ([]model.Result, error) {
	filtered := results
	if opts.Regex {
		pattern := query
		if opts.IgnoreCase {
			pattern = "(?i)" + pattern
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		filtered = filterByPredicate(filtered, func(item model.Result) bool {
			hay := item.Title + " " + strings.Join(item.Tags, " ")
			return re.MatchString(hay)
		})
	}

	if opts.Mood != "" {
		mood := opts.Mood
		if opts.IgnoreCase {
			mood = strings.ToLower(mood)
		}
		filtered = filterByPredicate(filtered, func(item model.Result) bool {
			hay := item.Title + " " + strings.Join(item.Tags, " ")
			if opts.IgnoreCase {
				hay = strings.ToLower(hay)
			}
			contains := strings.Contains(hay, mood)
			if opts.Invert {
				return !contains
			}
			return contains
		})
	}

	return filtered, nil
}

func filterByPredicate(items []model.Result, keep func(model.Result) bool) []model.Result {
	out := make([]model.Result, 0, len(items))
	for _, item := range items {
		if keep(item) {
			out = append(out, item)
		}
	}
	return out
}
