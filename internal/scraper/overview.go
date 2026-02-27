package scraper

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mathiswitte/flat-notifier-go/internal/model"
)

const baseURL = "https://www.kleinanzeigen.de"

// GetInfosFromOverview parses a Kleinanzeigen overview page and returns flat listings.
func GetInfosFromOverview(body io.Reader) ([]model.FlatInfo, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	var flats []model.FlatInfo

	doc.Find(".aditem-main").Each(func(_ int, s *goquery.Selection) {
		link := s.Find("a.ellipsis")
		href, exists := link.Attr("href")
		if !exists || href == "" {
			return
		}

		href = strings.TrimSpace(href)
		flatURL := href
		if strings.HasPrefix(href, "/") {
			flatURL = baseURL + href
		}

		flatID := GetIDFromFlatLink(href)
		if flatID == "" {
			return
		}
		flats = append(flats, model.FlatInfo{
			FlatURL: flatURL,
			FlatID:  flatID,
		})
	})

	return flats, nil
}
