package app

import (
	"context"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	"github.com/aws/aws-sdk-go-v2/service/ssm"

	dynamodbadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/dynamodb"
	envadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/env"
	postgresadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/postgres"
	ssmadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/ssm"
	"github.com/PlatformStackPulse/go-lambda-template/internal/config"
	"github.com/PlatformStackPulse/go-lambda-template/internal/domain"
	"github.com/PlatformStackPulse/go-lambda-template/internal/handler"
	"github.com/PlatformStackPulse/go-lambda-template/internal/logger"
	"github.com/PlatformStackPulse/go-lambda-template/internal/usecase"
)

type Application struct {
	Config  *config.Config
	Logger  *logger.Logger
	Handler *handler.APIHandler
}

func New(ctx context.Context) (*Application, error) {
	// Load runtime configuration from environment variables.
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Create a structured logger early so all downstream components can emit consistent logs.
	log := logger.NewLogger(cfg.Debug)

	// Build a shared AWS SDK configuration used by all AWS-backed adapters.
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.AWS.Region))
	if err != nil {
		return nil, err
	}

	// Read app-level configuration values from a centralized SSM JSON document.
	ssmProvider := ssmadapter.NewParameterStore(
		ssm.NewFromConfig(awsCfg),
		cfg.SSM.AppConfigParameterName,
	)

	// DynamoDB is always enabled as the default persistence adapter.
	recorder := dynamodbadapter.NewGreetingRecorder(
		dynamodb.NewFromConfig(awsCfg),
		cfg.DynamoDB.RequestsTableName,
	)

	// Fan out writes across all configured persistence backends.
	recorders := []usecase.GreetingRecorder{recorder}

	// Postgres is optional and only appended when explicitly enabled.
	if cfg.Postgres.Enabled {
		recorders = append(recorders, postgresadapter.NewGreetingRecorder(
			rdsdata.NewFromConfig(awsCfg),
			cfg.Postgres.ClusterARN,
			cfg.Postgres.SecretARN,
			cfg.Postgres.DatabaseName,
		))
	}

	// Runtime settings exposes environment-variable overrides to use cases.
	environment := envadapter.NewRuntimeSettings()

	// Compose the application use case and HTTP handler.
	greetingUseCase := usecase.NewGreetingUseCase(log, ssmProvider, &fanoutGreetingRecorder{recorders: recorders}, environment)
	apiHandler := handler.NewAPIHandler(log, greetingUseCase)

	// Return the fully wired application container used by the Lambda entrypoint.
	return &Application{
		Config:  cfg,
		Logger:  log,
		Handler: apiHandler,
	}, nil
}

type fanoutGreetingRecorder struct {
	recorders []usecase.GreetingRecorder
}

func (r *fanoutGreetingRecorder) Record(ctx context.Context, record domain.GreetingRecord) error {
	// Execute recorders in order and fail fast on the first write error.
	for _, recorder := range r.recorders {
		if recorder == nil {
			continue
		}
		if err := recorder.Record(ctx, record); err != nil {
			return fmt.Errorf("fanout recorder failed: %w", err)
		}
	}

	return nil
}
