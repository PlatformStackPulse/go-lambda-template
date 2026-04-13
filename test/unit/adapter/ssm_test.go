package adapter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ssmadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/ssm"
	apperrors "github.com/PlatformStackPulse/go-lambda-template/internal/errors"
)

type stubGetParameterClient struct {
	output *ssm.GetParameterOutput
	err    error
}

func (s *stubGetParameterClient) GetParameter(_ context.Context, _ *ssm.GetParameterInput, _ ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.output, nil
}

func TestParameterStoreStringValue(t *testing.T) {
	// Verify happy path extraction from the JSON app-config document.
	store := ssmadapter.NewParameterStore(&stubGetParameterClient{output: &ssm.GetParameterOutput{Parameter: &ssmtypes.Parameter{Value: aws.String(`{"sample.greeting.prefix":"Hello"}`)}}}, "/app/dev/app-config")

	prefix, err := store.StringValue(context.Background(), "sample.greeting.prefix")
	require.NoError(t, err)
	assert.Equal(t, "Hello", prefix)
}

func TestParameterStoreReturnsErrorForEmptyDocument(t *testing.T) {
	store := ssmadapter.NewParameterStore(&stubGetParameterClient{output: &ssm.GetParameterOutput{Parameter: &ssmtypes.Parameter{Value: aws.String(" ")}}}, "/app/dev/app-config")

	_, err := store.StringValue(context.Background(), "sample.greeting.prefix")
	require.Error(t, err)
	assert.True(t, apperrors.IsCode(err, apperrors.ErrConfiguration))
}

func TestParameterStoreReturnsInvalidJSONError(t *testing.T) {
	// Invalid JSON in SSM should fail fast as a configuration error.
	store := ssmadapter.NewParameterStore(&stubGetParameterClient{output: &ssm.GetParameterOutput{Parameter: &ssmtypes.Parameter{Value: aws.String("not-json")}}}, "/app/dev/app-config")

	_, err := store.StringValue(context.Background(), "sample.greeting.prefix")
	require.Error(t, err)
	assert.True(t, apperrors.IsCode(err, apperrors.ErrConfiguration))
}

func TestParameterStoreReturnsReadError(t *testing.T) {
	store := ssmadapter.NewParameterStore(&stubGetParameterClient{err: errors.New("boom")}, "/app/dev/app-config")

	_, err := store.StringValue(context.Background(), "sample.greeting.prefix")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read SSM parameter")
}
