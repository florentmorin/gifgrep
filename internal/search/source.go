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
		return "tenor"
	}
	return source
}
