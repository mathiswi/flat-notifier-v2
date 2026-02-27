package format

import "github.com/mathiswitte/flat-notifier-go/internal/model"

// ForTelegram formats a flat for Telegram using Markdown (*bold*).
func ForTelegram(flat model.Flat) string {
	return flat.URL + "\n\n" +
		flat.Title + "\n\n" +
		"*Kalt*: " + flat.ColdRent + "\n" +
		"*Warm*: " + flat.WarmRent + "\n" +
		"*Nebenkosten*: " + flat.ExtraCosts + "\n" +
		"*Größe*: " + flat.Size + "\n" +
		"*Adresse*: " + flat.Address + "\n" +
		"*Zimmer*: " + flat.Rooms + "\n" +
		"*Typ*: " + flat.ApartmentType + "\n" +
		"*Verfügbar:* " + flat.AvailableFrom
}

// ForDiscord formats a flat for Discord using Markdown (**bold**).
func ForDiscord(flat model.Flat) string {
	return flat.URL + "\n\n" +
		flat.Title + "\n\n" +
		"**Kalt**: " + flat.ColdRent + "\n" +
		"**Warm**: " + flat.WarmRent + "\n" +
		"**Nebenkosten**: " + flat.ExtraCosts + "\n" +
		"**Größe**: " + flat.Size + "\n" +
		"**Adresse**: " + flat.Address + "\n" +
		"**Zimmer**: " + flat.Rooms + "\n" +
		"**Typ**: " + flat.ApartmentType + "\n" +
		"**Verfügbar:** " + flat.AvailableFrom
}
