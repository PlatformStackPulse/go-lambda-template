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

func TestGreetingParameterStoreGreetingPrefix(t *testing.T) {
	store := ssmadapter.NewGreetingParameterStore(&stubGetParameterClient{output: &ssm.GetParameterOutput{Parameter: &ssmtypes.Parameter{Value: aws.String("Hello")}}}, "/app/dev/greeting-prefix")

	prefix, err := store.GreetingPrefix(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Hello", prefix)
}

func TestGreetingParameterStoreReturnsErrorForEmptyValue(t *testing.T) {
	store := ssmadapter.NewGreetingParameterStore(&stubGetParameterClient{output: &ssm.GetParameterOutput{Parameter: &ssmtypes.Parameter{Value: aws.String(" ")}}}, "/app/dev/greeting-prefix")

	_, err := store.GreetingPrefix(context.Background())
	require.Error(t, err)
	assert.True(t, apperrors.IsCode(err, apperrors.ErrConfiguration))
}

func TestGreetingParameterStoreReturnsReadError(t *testing.T) {
	store := ssmadapter.NewGreetingParameterStore(&stubGetParameterClient{err: errors.New("boom")}, "/app/dev/greeting-prefix")

	_, err := store.GreetingPrefix(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read SSM parameter")
}
