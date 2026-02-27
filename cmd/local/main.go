package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/mathiswitte/flat-notifier-go/internal/compare"
	"github.com/mathiswitte/flat-notifier-go/internal/db"
	"github.com/mathiswitte/flat-notifier-go/internal/format"
	"github.com/mathiswitte/flat-notifier-go/internal/notify"
	"github.com/mathiswitte/flat-notifier-go/internal/scraper"
)

func main() {
	ctx := context.Background()

	overviewURL := os.Getenv("OVERVIEW_URL")
	if overviewURL == "" {
		overviewURL = "https://www.kleinanzeigen.de/s-wohnung-mieten/26127/c203l3112"
	}

	// Connect to DynamoDB Local
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("eu-central-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("local", "local", "")),
	)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	// Ensure table exists
	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		tableName = "ebayTable"
	}
	ensureTable(ctx, client, tableName)

	store := db.NewDynamoStore(client)

	// Fetch overview
	body, err := scraper.FetchOverviewPage(overviewURL)
	if err != nil {
		log.Fatalf("failed to fetch overview: %v", err)
	}
	defer body.Close()

	flats, err := scraper.GetInfosFromOverview(body)
	if err != nil {
		log.Fatalf("failed to parse overview: %v", err)
	}
	fmt.Printf("Found %d flats in overview\n", len(flats))

	// Compare against stored IDs
	existingIDs, err := store.GetAllIDs(ctx)
	if err != nil {
		log.Fatalf("failed to get existing IDs: %v", err)
	}
	fmt.Printf("Found %d existing IDs in local DB\n", len(existingIDs))

	newFlats := compare.FindNewFlats(flats, existingIDs)
	if len(newFlats) == 0 {
		fmt.Println("No new flats found")
		return
	}
	fmt.Printf("Found %d new flats\n\n", len(newFlats))

	// Process each new flat
	dispatcher, _ := notify.NewDispatcher()

	for _, flat := range newFlats {
		// Store ID
		if err := store.WriteID(ctx, flat.FlatID); err != nil {
			log.Printf("error writing ID %s: %v", flat.FlatID, err)
			continue
		}

		// Scrape details
		details, err := scraper.FetchAndScrapeFlatPage(flat.FlatURL)
		if err != nil {
			log.Printf("error scraping %s: %v", flat.FlatURL, err)
			continue
		}

		fmt.Println("=== Discord Format ===")
		fmt.Println(format.ForDiscord(details))
		fmt.Println()

		// Send notification if configured
		if dispatcher != nil {
			msg := format.ForDiscord(details)
			if err := dispatcher.Send(msg, details.FirstImageURL); err != nil {
				log.Printf("notification error: %v", err)
			} else {
				fmt.Println("Notification sent!")
			}
		}
	}
}

func ensureTable(ctx context.Context, client *dynamodb.Client, tableName string) {
	_, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err == nil {
		return
	}

	_, err = client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []dtypes.KeySchemaElement{
			{AttributeName: aws.String("flatId"), KeyType: dtypes.KeyTypeHash},
		},
		AttributeDefinitions: []dtypes.AttributeDefinition{
			{AttributeName: aws.String("flatId"), AttributeType: dtypes.ScalarAttributeTypeS},
		},
		BillingMode: dtypes.BillingModePayPerRequest,
	})
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	log.Println("Created local DynamoDB table:", tableName)
}
