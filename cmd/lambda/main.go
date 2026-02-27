package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"golang.org/x/sync/errgroup"

	"github.com/mathiswitte/flat-notifier-go/internal/compare"
	"github.com/mathiswitte/flat-notifier-go/internal/db"
	"github.com/mathiswitte/flat-notifier-go/internal/format"
	"github.com/mathiswitte/flat-notifier-go/internal/model"
	"github.com/mathiswitte/flat-notifier-go/internal/notify"
	"github.com/mathiswitte/flat-notifier-go/internal/scraper"
)

func handler(ctx context.Context) error {
	overviewURL := os.Getenv("OVERVIEW_URL")
	if overviewURL == "" {
		return fmt.Errorf("OVERVIEW_URL environment variable not set")
	}

	// Initialize DynamoDB store
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-central-1"))
	if err != nil {
		return fmt.Errorf("loading AWS config: %w", err)
	}
	store := db.NewDynamoStore(dynamodb.NewFromConfig(cfg))

	// Initialize notification dispatcher
	dispatcher, err := notify.NewDispatcher()
	if err != nil {
		return fmt.Errorf("initializing notifications: %w", err)
	}

	// Step 1+2: Fetch overview + get existing IDs concurrently
	var foundFlats []model.FlatInfo
	var existingIDs []string

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		body, err := scraper.FetchOverviewPage(overviewURL)
		if err != nil {
			return fmt.Errorf("fetching overview: %w", err)
		}
		defer body.Close()

		foundFlats, err = scraper.GetInfosFromOverview(body)
		if err != nil {
			return fmt.Errorf("parsing overview: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		existingIDs, err = store.GetAllIDs(gctx)
		if err != nil {
			return fmt.Errorf("scanning DynamoDB: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	// Step 3: Compare for new flats
	newFlats := compare.FindNewFlats(foundFlats, existingIDs)
	if len(newFlats) == 0 {
		log.Println("Keine neue Wohnung verfügbar")
		return nil
	}

	log.Printf("Wohnung gefunden: %d neue", len(newFlats))

	// Step 4: For each new flat - write ID, scrape details, format, notify
	g2, _ := errgroup.WithContext(ctx)

	for _, flat := range newFlats {
		flat := flat // capture loop var
		g2.Go(func() error {
			// Write ID first (crash safety)
			if err := store.WriteID(ctx, flat.FlatID); err != nil {
				log.Printf("error writing ID %s: %v", flat.FlatID, err)
				return err
			}

			// Scrape detail page
			details, err := scraper.FetchAndScrapeFlatPage(flat.FlatURL)
			if err != nil {
				log.Printf("error scraping %s: %v", flat.FlatURL, err)
				return nil // Don't fail entire batch for one scrape error
			}

			// Format and send
			message := format.ForTelegram(details)
			if err := dispatcher.Send(message, details.FirstImageURL); err != nil {
				log.Printf("error sending notification for %s: %v", flat.FlatID, err)
			}
			return nil
		})
	}

	return g2.Wait()
}

func main() {
	lambda.Start(handler)
}
