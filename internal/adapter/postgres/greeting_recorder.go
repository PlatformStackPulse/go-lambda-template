package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	rdsdatatypes "github.com/aws/aws-sdk-go-v2/service/rdsdata/types"

	"github.com/PlatformStackPulse/go-lambda-template/internal/domain"
	apperrors "github.com/PlatformStackPulse/go-lambda-template/internal/errors"
)

const createGreetingRecordsTableSQL = `
CREATE TABLE IF NOT EXISTS greeting_records (
    request_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    source TEXT NOT NULL
)`

const insertGreetingRecordSQL = `
INSERT INTO greeting_records (request_id, name, message, created_at, source)
VALUES (:request_id, :name, :message, :created_at, :source)
ON CONFLICT (request_id) DO NOTHING`

const selectGreetingRecordByRequestIDSQL = `
SELECT request_id, name, message, created_at, source
FROM greeting_records
WHERE request_id = :request_id
LIMIT 1`

const selectGreetingRecordsSQL = `
SELECT request_id, name, message, created_at, source
FROM greeting_records
ORDER BY created_at DESC
LIMIT :limit`

type ExecuteStatementAPI interface {
	ExecuteStatement(context.Context, *rdsdata.ExecuteStatementInput, ...func(*rdsdata.Options)) (*rdsdata.ExecuteStatementOutput, error)
}

type GreetingRecorder struct {
	client      ExecuteStatementAPI
	resourceARN string
	secretARN   string
	database    string
}

func NewGreetingRecorder(client ExecuteStatementAPI, resourceARN, secretARN, database string) *GreetingRecorder {
	return &GreetingRecorder{client: client, resourceARN: resourceARN, secretARN: secretARN, database: database}
}

func (r *GreetingRecorder) Record(ctx context.Context, record domain.GreetingRecord) error {
	if err := r.execute(ctx, createGreetingRecordsTableSQL, nil); err != nil {
		return fmt.Errorf("ensure greeting_records table: %w", err)
	}

	params := []rdsdatatypes.SqlParameter{
		stringParam("request_id", record.RequestID),
		stringParam("name", record.Name),
		stringParam("message", record.Message),
		stringParam("created_at", record.CreatedAt),
		stringParam("source", record.Source),
	}

	if err := r.execute(ctx, insertGreetingRecordSQL, params); err != nil {
		return fmt.Errorf("insert greeting record: %w", err)
	}

	return nil
}

func (r *GreetingRecorder) GetByRequestID(ctx context.Context, requestID string) (domain.GreetingRecord, error) {
	if err := r.execute(ctx, createGreetingRecordsTableSQL, nil); err != nil {
		return domain.GreetingRecord{}, fmt.Errorf("ensure greeting_records table: %w", err)
	}

	output, err := r.executeStatement(ctx, selectGreetingRecordByRequestIDSQL, []rdsdatatypes.SqlParameter{
		stringParam("request_id", requestID),
	})
	if err != nil {
		return domain.GreetingRecord{}, fmt.Errorf("select greeting record by request id: %w", err)
	}
	if len(output.Records) == 0 {
		return domain.GreetingRecord{}, apperrors.New(apperrors.ErrNotFound, "greeting record not found")
	}

	record, err := greetingRecordFromFields(output.Records[0])
	if err != nil {
		return domain.GreetingRecord{}, fmt.Errorf("decode greeting record: %w", err)
	}

	return record, nil
}

func (r *GreetingRecorder) List(ctx context.Context, limit int) ([]domain.GreetingRecord, error) {
	if err := r.execute(ctx, createGreetingRecordsTableSQL, nil); err != nil {
		return nil, fmt.Errorf("ensure greeting_records table: %w", err)
	}

	if limit <= 0 {
		limit = 10
	}

	output, err := r.executeStatement(ctx, selectGreetingRecordsSQL, []rdsdatatypes.SqlParameter{
		longParam("limit", int64(limit)),
	})
	if err != nil {
		return nil, fmt.Errorf("select greeting records: %w", err)
	}

	records := make([]domain.GreetingRecord, 0, len(output.Records))
	for _, row := range output.Records {
		record, decodeErr := greetingRecordFromFields(row)
		if decodeErr != nil {
			return nil, fmt.Errorf("decode greeting record: %w", decodeErr)
		}
		records = append(records, record)
	}

	return records, nil
}

func (r *GreetingRecorder) execute(ctx context.Context, sql string, params []rdsdatatypes.SqlParameter) error {
	_, err := r.executeStatement(ctx, sql, params)
	return err
}

func (r *GreetingRecorder) executeStatement(ctx context.Context, sql string, params []rdsdatatypes.SqlParameter) (*rdsdata.ExecuteStatementOutput, error) {
	output, err := r.client.ExecuteStatement(ctx, &rdsdata.ExecuteStatementInput{
		ResourceArn: aws.String(r.resourceARN),
		SecretArn:   aws.String(r.secretARN),
		Database:    aws.String(r.database),
		Sql:         aws.String(sql),
		Parameters:  params,
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

func stringParam(name, value string) rdsdatatypes.SqlParameter {
	return rdsdatatypes.SqlParameter{
		Name: aws.String(name),
		Value: &rdsdatatypes.FieldMemberStringValue{
			Value: value,
		},
	}
}

func longParam(name string, value int64) rdsdatatypes.SqlParameter {
	return rdsdatatypes.SqlParameter{
		Name: aws.String(name),
		Value: &rdsdatatypes.FieldMemberLongValue{
			Value: value,
		},
	}
}

func greetingRecordFromFields(fields []rdsdatatypes.Field) (domain.GreetingRecord, error) {
	if len(fields) < 5 {
		return domain.GreetingRecord{}, fmt.Errorf("expected 5 columns, got %d", len(fields))
	}

	requestID, err := fieldString(fields[0])
	if err != nil {
		return domain.GreetingRecord{}, fmt.Errorf("request_id: %w", err)
	}
	name, err := fieldString(fields[1])
	if err != nil {
		return domain.GreetingRecord{}, fmt.Errorf("name: %w", err)
	}
	message, err := fieldString(fields[2])
	if err != nil {
		return domain.GreetingRecord{}, fmt.Errorf("message: %w", err)
	}
	createdAt, err := fieldString(fields[3])
	if err != nil {
		return domain.GreetingRecord{}, fmt.Errorf("created_at: %w", err)
	}
	source, err := fieldString(fields[4])
	if err != nil {
		return domain.GreetingRecord{}, fmt.Errorf("source: %w", err)
	}

	return domain.GreetingRecord{
		RequestID: requestID,
		Name:      name,
		Message:   message,
		CreatedAt: createdAt,
		Source:    source,
	}, nil
}

func fieldString(field rdsdatatypes.Field) (string, error) {
	switch value := field.(type) {
	case *rdsdatatypes.FieldMemberStringValue:
		return value.Value, nil
	case *rdsdatatypes.FieldMemberLongValue:
		return strconv.FormatInt(value.Value, 10), nil
	case *rdsdatatypes.FieldMemberDoubleValue:
		return strconv.FormatFloat(value.Value, 'f', -1, 64), nil
	case *rdsdatatypes.FieldMemberBooleanValue:
		return strconv.FormatBool(value.Value), nil
	case *rdsdatatypes.FieldMemberIsNull:
		if value.Value {
			return "", nil
		}
	}

	return "", fmt.Errorf("unsupported field type %T", field)
}
