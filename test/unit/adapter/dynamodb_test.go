package adapter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dynamodbadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/dynamodb"
	"github.com/PlatformStackPulse/go-lambda-template/internal/domain"
)

type stubPutItemClient struct {
	input *dynamodb.PutItemInput
	err   error
}

func (s *stubPutItemClient) PutItem(_ context.Context, input *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	s.input = input
	if s.err != nil {
		return nil, s.err
	}
	return &dynamodb.PutItemOutput{}, nil
}

func TestGreetingRecorderRecord(t *testing.T) {
	client := &stubPutItemClient{}
	recorder := dynamodbadapter.NewGreetingRecorder(client, "requests")

	err := recorder.Record(context.Background(), domain.GreetingRecord{RequestID: "req-1", Name: "Dev", Message: "Hello, Dev!"})
	require.NoError(t, err)
	require.NotNil(t, client.input)
	assert.Equal(t, "requests", *client.input.TableName)
	assert.Contains(t, client.input.Item, "request_id")
}

func TestGreetingRecorderReturnsWrappedError(t *testing.T) {
	recorder := dynamodbadapter.NewGreetingRecorder(&stubPutItemClient{err: errors.New("put failed")}, "requests")

	err := recorder.Record(context.Background(), domain.GreetingRecord{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "put greeting record")
}
