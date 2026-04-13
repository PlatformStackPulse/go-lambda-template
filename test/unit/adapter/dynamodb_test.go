package adapter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dynamodbadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/dynamodb"
	"github.com/PlatformStackPulse/go-lambda-template/internal/domain"
	apperrors "github.com/PlatformStackPulse/go-lambda-template/internal/errors"
)

type stubPutItemClient struct {
	putInput  *dynamodb.PutItemInput
	getInput  *dynamodb.GetItemInput
	scanInput *dynamodb.ScanInput
	getItem   map[string]types.AttributeValue
	scanItems []map[string]types.AttributeValue
	err       error
}

func (s *stubPutItemClient) PutItem(_ context.Context, input *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	s.putInput = input
	if s.err != nil {
		return nil, s.err
	}
	return &dynamodb.PutItemOutput{}, nil
}

func (s *stubPutItemClient) GetItem(_ context.Context, input *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	s.getInput = input
	if s.err != nil {
		return nil, s.err
	}
	return &dynamodb.GetItemOutput{Item: s.getItem}, nil
}

func (s *stubPutItemClient) Scan(_ context.Context, input *dynamodb.ScanInput, _ ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	s.scanInput = input
	if s.err != nil {
		return nil, s.err
	}
	return &dynamodb.ScanOutput{Items: s.scanItems}, nil
}

func TestGreetingRecorderRecord(t *testing.T) {
	client := &stubPutItemClient{}
	recorder := dynamodbadapter.NewGreetingRecorder(client, "requests")

	err := recorder.Record(context.Background(), domain.GreetingRecord{RequestID: "req-1", Name: "Dev", Message: "Hello, Dev!"})
	require.NoError(t, err)
	require.NotNil(t, client.putInput)
	assert.Equal(t, "requests", *client.putInput.TableName)
	assert.Contains(t, client.putInput.Item, "request_id")
}

func TestGreetingRecorderReturnsWrappedError(t *testing.T) {
	recorder := dynamodbadapter.NewGreetingRecorder(&stubPutItemClient{err: errors.New("put failed")}, "requests")

	err := recorder.Record(context.Background(), domain.GreetingRecord{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "put greeting record")
}

func TestGreetingRecorderGetByRequestID(t *testing.T) {
	item, err := attributevalue.MarshalMap(domain.GreetingRecord{
		RequestID: "req-1",
		Name:      "Dev",
		Message:   "Hello, Dev!",
		CreatedAt: "2026-04-13T12:00:00Z",
		Source:    "unit-test",
	})
	require.NoError(t, err)

	client := &stubPutItemClient{getItem: item}
	recorder := dynamodbadapter.NewGreetingRecorder(client, "requests")

	record, err := recorder.GetByRequestID(context.Background(), "req-1")
	require.NoError(t, err)
	assert.Equal(t, "req-1", record.RequestID)
	require.NotNil(t, client.getInput)
	assert.Equal(t, "requests", *client.getInput.TableName)
}

func TestGreetingRecorderGetByRequestIDReturnsNotFound(t *testing.T) {
	recorder := dynamodbadapter.NewGreetingRecorder(&stubPutItemClient{}, "requests")

	_, err := recorder.GetByRequestID(context.Background(), "missing")
	require.Error(t, err)
	assert.True(t, apperrors.IsCode(err, apperrors.ErrNotFound))
}

func TestGreetingRecorderList(t *testing.T) {
	first, err := attributevalue.MarshalMap(domain.GreetingRecord{
		RequestID: "req-1",
		Name:      "Ada",
		Message:   "Hello, Ada!",
		CreatedAt: "2026-04-13T12:00:00Z",
		Source:    "unit-test",
	})
	require.NoError(t, err)
	second, err := attributevalue.MarshalMap(domain.GreetingRecord{
		RequestID: "req-2",
		Name:      "Grace",
		Message:   "Hello, Grace!",
		CreatedAt: "2026-04-13T12:01:00Z",
		Source:    "unit-test",
	})
	require.NoError(t, err)

	client := &stubPutItemClient{scanItems: []map[string]types.AttributeValue{first, second}}
	recorder := dynamodbadapter.NewGreetingRecorder(client, "requests")

	records, err := recorder.List(context.Background(), 2)
	require.NoError(t, err)
	assert.Len(t, records, 2)
	assert.Equal(t, "req-2", records[0].RequestID)
	assert.Equal(t, int32(2), *client.scanInput.Limit)
}

func TestGreetingRecorderListReturnsWrappedError(t *testing.T) {
	recorder := dynamodbadapter.NewGreetingRecorder(&stubPutItemClient{err: errors.New("scan failed")}, "requests")

	_, err := recorder.List(context.Background(), 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scan greeting records")
}
