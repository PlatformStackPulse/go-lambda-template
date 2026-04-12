package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/PlatformStackPulse/go-lambda-template/internal/domain"
)

type PutItemAPI interface {
	PutItem(context.Context, *awsdynamodb.PutItemInput, ...func(*awsdynamodb.Options)) (*awsdynamodb.PutItemOutput, error)
}

type GreetingRecorder struct {
	client    PutItemAPI
	tableName string
}

func NewGreetingRecorder(client PutItemAPI, tableName string) *GreetingRecorder {
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