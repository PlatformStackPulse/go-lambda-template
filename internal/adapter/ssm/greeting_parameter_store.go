package ssm

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsssm "github.com/aws/aws-sdk-go-v2/service/ssm"

	apperrors "github.com/PlatformStackPulse/go-lambda-template/internal/errors"
)

type GetParameterAPI interface {
	GetParameter(context.Context, *awsssm.GetParameterInput, ...func(*awsssm.Options)) (*awsssm.GetParameterOutput, error)
}

type GreetingParameterStore struct {
	client        GetParameterAPI
	parameterName string
}

func NewGreetingParameterStore(client GetParameterAPI, parameterName string) *GreetingParameterStore {
	return &GreetingParameterStore{client: client, parameterName: parameterName}
}

func (s *GreetingParameterStore) GreetingPrefix(ctx context.Context) (string, error) {
	output, err := s.client.GetParameter(ctx, &awsssm.GetParameterInput{
		Name:           aws.String(s.parameterName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("read SSM parameter %q: %w", s.parameterName, err)
	}

	if output.Parameter == nil || strings.TrimSpace(aws.ToString(output.Parameter.Value)) == "" {
		return "", apperrors.New(apperrors.ErrConfiguration, "greeting parameter value is empty")
	}

	return strings.TrimSpace(aws.ToString(output.Parameter.Value)), nil
}