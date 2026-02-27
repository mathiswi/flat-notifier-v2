package scraper

import (
	"os"
	"strings"
	"testing"
)

func TestScrapeFlatPage(t *testing.T) {
	f, err := os.Open("../../testdata/detail.html")
	if err != nil {
		t.Fatalf("failed to open testdata: %v", err)
	}
	defer f.Close()

	flat, err := ScrapeFlatPage(f)
	if err != nil {
		t.Fatalf("ScrapeFlatPage returned error: %v", err)
	}

	if flat.Title == "" {
		t.Error("Title is empty")
	}
	if flat.Address == "" {
		t.Error("Address is empty")
	}
	if flat.ColdRent == "" {
		t.Error("ColdRent is empty")
	}

	t.Logf("Title: %q", flat.Title)
	t.Logf("Address: %q", flat.Address)
	t.Logf("ColdRent: %q", flat.ColdRent)
	t.Logf("WarmRent: %q", flat.WarmRent)
	t.Logf("ExtraCosts: %q", flat.ExtraCosts)
	t.Logf("Size: %q", flat.Size)
	t.Logf("Rooms: %q", flat.Rooms)
	t.Logf("ApartmentType: %q", flat.ApartmentType)
	t.Logf("AvailableFrom: %q", flat.AvailableFrom)
	t.Logf("Features: %v", flat.Features)
	t.Logf("FirstImageURL: %q", flat.FirstImageURL)

	// Verify specific values from the test HTML (3337672397 detail page)
	if flat.Title != "Helle 2 Zimmer Wohnung inkl. Küche in Ofenerdiek zu vermieten" {
		t.Errorf("Title = %q, want exact match", flat.Title)
	}
	if flat.ColdRent != "415 €" {
		t.Errorf("ColdRent = %q, want %q", flat.ColdRent, "415 €")
	}
	if flat.Size != "45 m²" {
		t.Errorf("Size = %q, want %q", flat.Size, "45 m²")
	}
	if flat.Rooms != "2" {
		t.Errorf("Rooms = %q, want %q", flat.Rooms, "2")
	}
	if flat.ApartmentType != "Etagenwohnung" {
		t.Errorf("ApartmentType = %q, want %q", flat.ApartmentType, "Etagenwohnung")
	}
	if flat.AvailableFrom != "März 2026" {
		t.Errorf("AvailableFrom = %q, want %q", flat.AvailableFrom, "März 2026")
	}
	if flat.ExtraCosts != "130 €" {
		t.Errorf("ExtraCosts = %q, want %q", flat.ExtraCosts, "130 €")
	}
	if flat.FirstImageURL == "" {
		t.Error("FirstImageURL is empty")
	}
	if len(flat.Features) != 4 {
		t.Errorf("Features count = %d, want 4", len(flat.Features))
	}
}

func TestScrapeFlatPage_EmptyHTML(t *testing.T) {
	flat, err := ScrapeFlatPage(strings.NewReader("<html><body></body></html>"))
	if err != nil {
		t.Fatalf("ScrapeFlatPage on empty HTML returned error: %v", err)
	}

	if flat.Title != "" {
		t.Errorf("Title = %q, want empty", flat.Title)
	}
	if flat.Address != "" {
		t.Errorf("Address = %q, want empty", flat.Address)
	}
	if flat.ColdRent != "" {
		t.Errorf("ColdRent = %q, want empty", flat.ColdRent)
	}
	if len(flat.Features) != 0 {
		t.Errorf("Features = %v, want empty", flat.Features)
	}
	if flat.FirstImageURL != "" {
		t.Errorf("FirstImageURL = %q, want empty", flat.FirstImageURL)
	}
}
