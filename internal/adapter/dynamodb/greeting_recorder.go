package dynamodb

import (
	"context"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/PlatformStackPulse/go-lambda-template/internal/domain"
	apperrors "github.com/PlatformStackPulse/go-lambda-template/internal/errors"
)

type DynamoDBAPI interface {
	PutItem(context.Context, *awsdynamodb.PutItemInput, ...func(*awsdynamodb.Options)) (*awsdynamodb.PutItemOutput, error)
	GetItem(context.Context, *awsdynamodb.GetItemInput, ...func(*awsdynamodb.Options)) (*awsdynamodb.GetItemOutput, error)
	Scan(context.Context, *awsdynamodb.ScanInput, ...func(*awsdynamodb.Options)) (*awsdynamodb.ScanOutput, error)
}

type GreetingRecorder struct {
	client    DynamoDBAPI
	tableName string
}

func NewGreetingRecorder(client DynamoDBAPI, tableName string) *GreetingRecorder {
	return &GreetingRecorder{client: client, tableName: tableName}
}

func (r *GreetingRecorder) Record(ctx context.Context, record domain.GreetingRecord) error {
	item, err := attributevalue.MarshalMap(record)
	if err != nil {
		return fmt.Errorf("marshal greeting record: %w", err)
	}

	_, err = r.client.PutItem(ctx, &awsdynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put greeting record: %w", err)
	}

	return nil
}

func (r *GreetingRecorder) GetByRequestID(ctx context.Context, requestID string) (domain.GreetingRecord, error) {
	output, err := r.client.GetItem(ctx, &awsdynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"request_id": &types.AttributeValueMemberS{Value: requestID},
		},
	})
	if err != nil {
		return domain.GreetingRecord{}, fmt.Errorf("get greeting record: %w", err)
	}
	if len(output.Item) == 0 {
		return domain.GreetingRecord{}, apperrors.New(apperrors.ErrNotFound, "greeting record not found")
	}

	var record domain.GreetingRecord
	if err := attributevalue.UnmarshalMap(output.Item, &record); err != nil {
		return domain.GreetingRecord{}, fmt.Errorf("unmarshal greeting record: %w", err)
	}

	return record, nil
}

func (r *GreetingRecorder) List(ctx context.Context, limit int32) ([]domain.GreetingRecord, error) {
	if limit <= 0 {
		limit = 10
	}

	output, err := r.client.Scan(ctx, &awsdynamodb.ScanInput{
		TableName: aws.String(r.tableName),
		Limit:     aws.Int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("scan greeting records: %w", err)
	}

	var records []domain.GreetingRecord
	if err := attributevalue.UnmarshalListOfMaps(output.Items, &records); err != nil {
		return nil, fmt.Errorf("unmarshal greeting records: %w", err)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt > records[j].CreatedAt
	})

	return records, nil
}
