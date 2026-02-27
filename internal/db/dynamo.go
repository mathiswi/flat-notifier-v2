package db

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// FlatStore defines the interface for flat ID persistence.
type FlatStore interface {
	GetAllIDs(ctx context.Context) ([]string, error)
	WriteID(ctx context.Context, id string) error
}

// DynamoStore implements FlatStore using DynamoDB.
type DynamoStore struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoStore creates a new DynamoDB-backed flat store.
func NewDynamoStore(client *dynamodb.Client) *DynamoStore {
	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		tableName = "ebayTable"
	}
	return &DynamoStore{
		client:    client,
		tableName: tableName,
	}
}

// GetAllIDs retrieves all flat IDs from DynamoDB using paginated Scan.
func (s *DynamoStore) GetAllIDs(ctx context.Context) ([]string, error) {
	var ids []string
	var exclusiveStartKey map[string]types.AttributeValue

	for {
		input := &dynamodb.ScanInput{
			TableName:         aws.String(s.tableName),
			ExclusiveStartKey: exclusiveStartKey,
		}

		result, err := s.client.Scan(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, item := range result.Items {
			if v, ok := item["flatId"]; ok {
				if sv, ok := v.(*types.AttributeValueMemberS); ok {
					ids = append(ids, sv.Value)
				}
			}
		}

		if result.LastEvaluatedKey == nil {
			break
		}
		exclusiveStartKey = result.LastEvaluatedKey
	}

	return ids, nil
}

// WriteID stores a single flat ID in DynamoDB.
func (s *DynamoStore) WriteID(ctx context.Context, id string) error {
	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item: map[string]types.AttributeValue{
			"flatId": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}
