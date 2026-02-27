package format

import (
	"strings"
	"testing"

	"github.com/mathiswitte/flat-notifier-go/internal/model"
)

var testFlat = model.Flat{
	URL:           "https://www.kleinanzeigen.de/s-anzeige/test/123-203-3124",
	Title:         "Schöne 3-Zimmer Wohnung",
	Address:       "49080 Niedersachsen - Osnabrück",
	ColdRent:      "690 €",
	WarmRent:      "940 €",
	ExtraCosts:    "75 €",
	Size:          "65 m²",
	Rooms:         "3",
	ApartmentType: "Etagenwohnung",
	AvailableFrom: "Februar 2026",
}

func TestForTelegram(t *testing.T) {
	msg := ForTelegram(testFlat)

	checks := []string{
		testFlat.URL,
		testFlat.Title,
		"*Kalt*: 690 €",
		"*Warm*: 940 €",
		"*Nebenkosten*: 75 €",
		"*Größe*: 65 m²",
		"*Adresse*: 49080",
		"*Zimmer*: 3",
		"*Typ*: Etagenwohnung",
		"*Verfügbar:* Februar 2026",
	}

	for _, check := range checks {
		if !strings.Contains(msg, check) {
			t.Errorf("message missing %q", check)
		}
	}
}

func TestForDiscord(t *testing.T) {
	msg := ForDiscord(testFlat)

	checks := []string{
		testFlat.URL,
		testFlat.Title,
		"**Kalt**: 690 €",
		"**Warm**: 940 €",
	}

	for _, check := range checks {
		if !strings.Contains(msg, check) {
			t.Errorf("message missing %q", check)
		}
	}
}
