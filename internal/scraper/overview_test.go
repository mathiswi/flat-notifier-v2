package scraper

import (
	"os"
	"strings"
	"testing"
)

func TestGetInfosFromOverview(t *testing.T) {
	f, err := os.Open("../../testdata/overview.html")
	if err != nil {
		t.Fatalf("failed to open testdata: %v", err)
	}
	defer f.Close()

	flats, err := GetInfosFromOverview(f)
	if err != nil {
		t.Fatalf("GetInfosFromOverview returned error: %v", err)
	}

	if len(flats) == 0 {
		t.Fatal("expected at least one flat, got 0")
	}

	for i, flat := range flats {
		if flat.FlatID == "" {
			t.Errorf("flat[%d]: FlatID is empty", i)
		}
		if flat.FlatURL == "" {
			t.Errorf("flat[%d]: FlatURL is empty", i)
		}
		if !strings.Contains(flat.FlatURL, "kleinanzeigen.de") {
			t.Errorf("flat[%d]: FlatURL %q does not contain base URL", i, flat.FlatURL)
		}
	}

	t.Logf("Found %d flats", len(flats))
	if len(flats) > 0 {
		t.Logf("First flat: ID=%s URL=%s", flats[0].FlatID, flats[0].FlatURL)
	}
}
