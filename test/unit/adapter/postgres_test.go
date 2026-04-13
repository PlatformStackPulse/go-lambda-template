package adapter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	rdsdatatypes "github.com/aws/aws-sdk-go-v2/service/rdsdata/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	postgresadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/postgres"
	"github.com/PlatformStackPulse/go-lambda-template/internal/domain"
	apperrors "github.com/PlatformStackPulse/go-lambda-template/internal/errors"
)

type stubExecuteStatementClient struct {
	inputs  []*rdsdata.ExecuteStatementInput
	outputs []*rdsdata.ExecuteStatementOutput
	errAt   int
	err     error
}

func (s *stubExecuteStatementClient) ExecuteStatement(_ context.Context, input *rdsdata.ExecuteStatementInput, _ ...func(*rdsdata.Options)) (*rdsdata.ExecuteStatementOutput, error) {
	s.inputs = append(s.inputs, input)
	if s.err != nil && len(s.inputs) == s.errAt {
		return nil, s.err
	}
	if len(s.outputs) >= len(s.inputs) {
		return s.outputs[len(s.inputs)-1], nil
	}
	return &rdsdata.ExecuteStatementOutput{}, nil
}

func TestPostgresGreetingRecorderRecord(t *testing.T) {
	client := &stubExecuteStatementClient{}
	recorder := postgresadapter.NewGreetingRecorder(client, "cluster-arn", "secret-arn", "app")

	err := recorder.Record(context.Background(), domain.GreetingRecord{
		RequestID: "req-1",
		Name:      "Ada",
		Message:   "Hello, Ada!",
		CreatedAt: "2026-04-13T12:00:00Z",
		Source:    "unit-test",
	})
	require.NoError(t, err)
	require.Len(t, client.inputs, 2)
	assert.Contains(t, aws.ToString(client.inputs[0].Sql), "CREATE TABLE IF NOT EXISTS greeting_records")
	assert.Contains(t, aws.ToString(client.inputs[1].Sql), "INSERT INTO greeting_records")
	assert.Equal(t, "cluster-arn", aws.ToString(client.inputs[1].ResourceArn))
	assert.Equal(t, "secret-arn", aws.ToString(client.inputs[1].SecretArn))
	assert.Equal(t, "app", aws.ToString(client.inputs[1].Database))
	assert.Len(t, client.inputs[1].Parameters, 5)
}

func TestPostgresGreetingRecorderReturnsTableError(t *testing.T) {
	recorder := postgresadapter.NewGreetingRecorder(&stubExecuteStatementClient{errAt: 1, err: errors.New("boom")}, "cluster-arn", "secret-arn", "app")

	err := recorder.Record(context.Background(), domain.GreetingRecord{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ensure greeting_records table")
}

func TestPostgresGreetingRecorderReturnsInsertError(t *testing.T) {
	recorder := postgresadapter.NewGreetingRecorder(&stubExecuteStatementClient{errAt: 2, err: errors.New("boom")}, "cluster-arn", "secret-arn", "app")

	err := recorder.Record(context.Background(), domain.GreetingRecord{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "insert greeting record")
}

func TestPostgresGreetingRecorderGetByRequestID(t *testing.T) {
	client := &stubExecuteStatementClient{outputs: []*rdsdata.ExecuteStatementOutput{
		{},
		{
			Records: [][]rdsdatatypes.Field{{
				&rdsdatatypes.FieldMemberStringValue{Value: "req-1"},
				&rdsdatatypes.FieldMemberStringValue{Value: "Ada"},
				&rdsdatatypes.FieldMemberStringValue{Value: "Hello, Ada!"},
				&rdsdatatypes.FieldMemberStringValue{Value: "2026-04-13T12:00:00Z"},
				&rdsdatatypes.FieldMemberStringValue{Value: "unit-test"},
			}},
		},
	}}
	recorder := postgresadapter.NewGreetingRecorder(client, "cluster-arn", "secret-arn", "app")

	record, err := recorder.GetByRequestID(context.Background(), "req-1")
	require.NoError(t, err)
	assert.Equal(t, "req-1", record.RequestID)
	assert.Equal(t, "Ada", record.Name)
	assert.Equal(t, "Hello, Ada!", record.Message)
	assert.Len(t, client.inputs, 2)
	assert.Contains(t, aws.ToString(client.inputs[1].Sql), "WHERE request_id = :request_id")
}

func TestPostgresGreetingRecorderGetByRequestIDReturnsNotFound(t *testing.T) {
	client := &stubExecuteStatementClient{outputs: []*rdsdata.ExecuteStatementOutput{{}, {}}}
	recorder := postgresadapter.NewGreetingRecorder(client, "cluster-arn", "secret-arn", "app")

	_, err := recorder.GetByRequestID(context.Background(), "missing")
	require.Error(t, err)
	assert.True(t, apperrors.IsCode(err, apperrors.ErrNotFound))
}

func TestPostgresGreetingRecorderList(t *testing.T) {
	client := &stubExecuteStatementClient{outputs: []*rdsdata.ExecuteStatementOutput{
		{},
		{
			Records: [][]rdsdatatypes.Field{
				{
					&rdsdatatypes.FieldMemberStringValue{Value: "req-2"},
					&rdsdatatypes.FieldMemberStringValue{Value: "Grace"},
					&rdsdatatypes.FieldMemberStringValue{Value: "Hello, Grace!"},
					&rdsdatatypes.FieldMemberStringValue{Value: "2026-04-13T12:01:00Z"},
					&rdsdatatypes.FieldMemberStringValue{Value: "unit-test"},
				},
				{
					&rdsdatatypes.FieldMemberStringValue{Value: "req-1"},
					&rdsdatatypes.FieldMemberStringValue{Value: "Ada"},
					&rdsdatatypes.FieldMemberStringValue{Value: "Hello, Ada!"},
					&rdsdatatypes.FieldMemberStringValue{Value: "2026-04-13T12:00:00Z"},
					&rdsdatatypes.FieldMemberStringValue{Value: "unit-test"},
				},
			},
		},
	}}
	recorder := postgresadapter.NewGreetingRecorder(client, "cluster-arn", "secret-arn", "app")

	records, err := recorder.List(context.Background(), 2)
	require.NoError(t, err)
	assert.Len(t, records, 2)
	assert.Equal(t, "req-2", records[0].RequestID)
	assert.Equal(t, "req-1", records[1].RequestID)
	assert.Contains(t, aws.ToString(client.inputs[1].Sql), "ORDER BY created_at DESC")
	assert.Len(t, client.inputs[1].Parameters, 1)
}
