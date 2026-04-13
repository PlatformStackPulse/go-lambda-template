package ssm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsssm "github.com/aws/aws-sdk-go-v2/service/ssm"

	apperrors "github.com/PlatformStackPulse/go-lambda-template/internal/errors"
)

type GetParameterAPI interface {
	GetParameter(context.Context, *awsssm.GetParameterInput, ...func(*awsssm.Options)) (*awsssm.GetParameterOutput, error)
}

type ParameterStore struct {
	client        GetParameterAPI
	parameterName string
}

func NewParameterStore(client GetParameterAPI, parameterName string) *ParameterStore {
	return &ParameterStore{client: client, parameterName: parameterName}
}

func (s *ParameterStore) StringValue(ctx context.Context, key string) (string, error) {
	// Read one JSON document from SSM and resolve keys from that shared config payload.
	output, err := s.client.GetParameter(ctx, &awsssm.GetParameterInput{
		Name:           aws.String(s.parameterName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("read SSM parameter %q: %w", s.parameterName, err)
	}

	if output.Parameter == nil || strings.TrimSpace(aws.ToString(output.Parameter.Value)) == "" {
		return "", apperrors.New(apperrors.ErrConfiguration, "app config parameter value is empty")
	}

	values := map[string]string{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(aws.ToString(output.Parameter.Value))), &values); err != nil {
		return "", apperrors.Wrap(apperrors.ErrConfiguration, "app config parameter contains invalid JSON", err)
	}

	// Missing or blank keys are treated as configuration errors to avoid silent defaults.
	value := strings.TrimSpace(values[key])
	if value == "" {
		return "", apperrors.New(apperrors.ErrConfiguration, fmt.Sprintf("app config value %q is empty or missing", key))
	}

	return value, nil
}
