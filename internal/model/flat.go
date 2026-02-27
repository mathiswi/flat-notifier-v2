package model

// FlatInfo holds the URL and ID of a flat listing from the overview page.
type FlatInfo struct {
	FlatURL string
	FlatID  string
}

// Flat holds all scraped details for a single apartment listing.
type Flat struct {
	URL            string
	Title          string
	Address        string
	ColdRent       string
	WarmRent       string
	ExtraCosts     string
	Size           string
	Rooms          string
	ApartmentType  string
	AvailableFrom  string
	Features       []string
	FirstImageURL  string
}
