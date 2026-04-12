package app

import (
	"context"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ssm"

	dynamodbadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/dynamodb"
	envadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/env"
	ssmadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/ssm"
	"github.com/PlatformStackPulse/go-lambda-template/internal/config"
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
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log := logger.NewLogger(cfg.Debug)

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.AWSRegion))
	if err != nil {
		return nil, err
	}

	ssmProvider := ssmadapter.NewGreetingParameterStore(
		ssm.NewFromConfig(awsCfg),
		cfg.GreetingParameterName,
	)
	recorder := dynamodbadapter.NewGreetingRecorder(
		dynamodb.NewFromConfig(awsCfg),
		cfg.DynamoDBTableName,
	)
	environment := envadapter.NewRuntimeSettings()
	greetingUseCase := usecase.NewGreetingUseCase(log, ssmProvider, recorder, environment)
	apiHandler := handler.NewAPIHandler(log, greetingUseCase)

	return &Application{
		Config:  cfg,
		Logger:  log,
		Handler: apiHandler,
	}, nil
}
