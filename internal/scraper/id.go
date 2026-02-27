package scraper

import "strings"

// GetIDFromFlatLink extracts the numeric flat ID from a Kleinanzeigen URL.
// e.g. "/s-anzeige/some-title/1955972992-203-3112" -> "1955972992"
func GetIDFromFlatLink(link string) string {
	link = strings.TrimRight(link, "/")
	if link == "" {
		return ""
	}
	parts := strings.Split(link, "/")
	lastPart := parts[len(parts)-1]
	id := strings.Split(lastPart, "-")[0]
	if id == "" {
		return ""
	}
	for _, c := range id {
		if c < '0' || c > '9' {
			return ""
		}
	}
	return id
}
