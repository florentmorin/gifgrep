package search

import (
	"os"
	"strings"
)

func ResolveSource(source string) string {
	source = strings.ToLower(strings.TrimSpace(source))
	if source == "" || source == "auto" {
		if os.Getenv("GIPHY_API_KEY") != "" {
			return "giphy"
		}
		if os.Getenv("HEYPSTER_API_KEY") != "" {
			return "heypster"
		}
		return "tenor"
	}
	return source
}
