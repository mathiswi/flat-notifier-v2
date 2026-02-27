package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mathiswitte/flat-notifier-go/internal/model"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

var httpClient = &http.Client{Timeout: 15 * time.Second}

// ScrapeFlatPage parses a Kleinanzeigen detail page and returns flat details.
func ScrapeFlatPage(body io.Reader) (model.Flat, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return model.Flat{}, err
	}

	flat := model.Flat{
		Title:    trimAll(doc.Find("#viewad-title").Text()),
		Address:  trimAll(doc.Find("#viewad-locality").First().Text()),
		ColdRent: trimAll(doc.Find("#viewad-price").Text()),
	}

	doc.Find("li.addetailslist--detail").Each(func(_ int, s *goquery.Selection) {
		// The label is the text node before the value span.
		// Get the full text, then subtract the value to get the label.
		valueHTML, _ := s.Find(".addetailslist--detail--value").Html()
		value := trimAll(valueHTML)

		// Clone and remove value span to get just the label text
		label := trimAll(s.Clone().Children().Remove().End().Text())

		switch label {
		case "Wohnfläche":
			flat.Size = value
		case "Zimmer":
			flat.Rooms = value
		case "Wohnungstyp":
			flat.ApartmentType = value
		case "Verfügbar ab":
			flat.AvailableFrom = value
		case "Warmmiete":
			flat.WarmRent = value
		case "Nebenkosten":
			flat.ExtraCosts = value
		}
	})

	doc.Find("li.checktag").Each(func(_ int, s *goquery.Selection) {
		feature := trimAll(s.Text())
		if feature != "" {
			flat.Features = append(flat.Features, feature)
		}
	})

	imgSrc, exists := doc.Find("#viewad-image").First().Attr("src")
	if exists {
		flat.FirstImageURL = strings.TrimSpace(imgSrc)
	}

	return flat, nil
}

// FetchAndScrapeFlatPage fetches a detail page URL and scrapes it.
func FetchAndScrapeFlatPage(url string) (model.Flat, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return model.Flat{}, fmt.Errorf("creating request for %s: %w", url, err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "de-DE,de;q=0.9,en;q=0.8")

	resp, err := httpClient.Do(req)
	if err != nil {
		return model.Flat{}, fmt.Errorf("fetching %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Flat{}, fmt.Errorf("fetching %s: status %d", url, resp.StatusCode)
	}

	flat, err := ScrapeFlatPage(resp.Body)
	if err != nil {
		return model.Flat{}, err
	}
	flat.URL = url
	return flat, nil
}

// FetchOverviewPage fetches the overview page HTML.
func FetchOverviewPage(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "de-DE,de;q=0.9,en;q=0.8")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching overview: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("fetching overview: status %d", resp.StatusCode)
	}
	return resp.Body, nil
}

func trimAll(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
