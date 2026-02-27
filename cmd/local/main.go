package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mathiswitte/flat-notifier-go/internal/format"
	"github.com/mathiswitte/flat-notifier-go/internal/notify"
	"github.com/mathiswitte/flat-notifier-go/internal/scraper"
)

func main() {
	// Parse overview from testdata
	overviewFile, err := os.Open("testdata/overview.html")
	if err != nil {
		log.Fatalf("failed to open overview testdata: %v", err)
	}
	defer overviewFile.Close()

	flats, err := scraper.GetInfosFromOverview(overviewFile)
	if err != nil {
		log.Fatalf("failed to parse overview: %v", err)
	}

	fmt.Printf("Found %d flats in overview\n\n", len(flats))

	// Parse detail from testdata
	detailFile, err := os.Open("testdata/detail.html")
	if err != nil {
		log.Fatalf("failed to open detail testdata: %v", err)
	}
	defer detailFile.Close()

	flat, err := scraper.ScrapeFlatPage(detailFile)
	if err != nil {
		log.Fatalf("failed to parse detail: %v", err)
	}
	flat.URL = "https://www.kleinanzeigen.de/s-anzeige/example/123-203-3124"

	// Print formatted messages
	fmt.Println("=== Telegram Format ===")
	fmt.Println(format.ForTelegram(flat))
	fmt.Println()
	fmt.Println("=== Discord Format ===")
	fmt.Println(format.ForDiscord(flat))

	// Optionally send real notifications if env vars are set
	if os.Getenv("TELEGRAM_TOKEN") != "" || os.Getenv("DISCORD_TOKEN") != "" {
		dispatcher, err := notify.NewDispatcher()
		if err != nil {
			log.Fatalf("failed to create dispatcher: %v", err)
		}

		msg := format.ForTelegram(flat)
		if err := dispatcher.Send(msg, flat.FirstImageURL); err != nil {
			log.Printf("notification error: %v", err)
		} else {
			fmt.Println("\nNotification sent successfully!")
		}
	}
}
