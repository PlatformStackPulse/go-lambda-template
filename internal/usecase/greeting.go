package usecase

import (
	"context"
	"time"

	"github.com/PlatformStackPulse/go-lambda-template/internal/domain"
	apperrors "github.com/PlatformStackPulse/go-lambda-template/internal/errors"
	"github.com/PlatformStackPulse/go-lambda-template/internal/logger"
)

type GreetingConfigProvider interface {
	StringValue(context.Context, string) (string, error)
}

type GreetingRecorder interface {
	Record(context.Context, domain.GreetingRecord) error
}

type GreetingEnvironment interface {
	Lookup(context.Context, string) string
}

type Clock func() time.Time

type GreetingInput struct {
	Name      string
	RequestID string
	Source    string
}

type GreetingOutput struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

type GreetingUseCase struct {
	log            *logger.Logger
	configProvider GreetingConfigProvider
	recorder       GreetingRecorder
	environment    GreetingEnvironment
	now            Clock
}

const (
	sampleGreetingPrefixKey    = "sample.greeting.prefix"
	sampleGreetingPrefixEnvKey = "SAMPLE_GREETING_PREFIX"
	apiSourceLabelEnvKey       = "API_SOURCE_LABEL"
)

func NewGreetingUseCase(log *logger.Logger, configProvider GreetingConfigProvider, recorder GreetingRecorder, environment GreetingEnvironment) *GreetingUseCase {
	return NewGreetingUseCaseWithClock(log, configProvider, recorder, environment, time.Now)
}

func NewGreetingUseCaseWithClock(log *logger.Logger, configProvider GreetingConfigProvider, recorder GreetingRecorder, environment GreetingEnvironment, now Clock) *GreetingUseCase {
	return &GreetingUseCase{
		log:            log,
		configProvider: configProvider,
		recorder:       recorder,
		environment:    environment,
		now:            now,
	}
}

func (uc *GreetingUseCase) Execute(ctx context.Context, input GreetingInput) (GreetingOutput, error) {
	if uc.configProvider == nil || uc.recorder == nil {
		return GreetingOutput{}, apperrors.New(apperrors.ErrConfiguration, "use case dependencies are not configured")
	}

	prefix := ""
	if uc.environment != nil {
		prefix = uc.environment.Lookup(ctx, sampleGreetingPrefixEnvKey)
	}
	if prefix == "" {
		loadedPrefix, err := uc.configProvider.StringValue(ctx, sampleGreetingPrefixKey)
		if err != nil {
			return GreetingOutput{}, apperrors.Wrap(apperrors.ErrIntegration, "failed to load greeting prefix", err)
		}
		prefix = loadedPrefix
	}

	requestID := input.RequestID
	if requestID == "" {
		requestID = "unknown"
	}

	source := input.Source
	if source == "" {
		source = "api"
	}
	if uc.environment != nil {
		if sourceLabel := uc.environment.Lookup(ctx, apiSourceLabelEnvKey); sourceLabel != "" {
			source = sourceLabel
		}
	}

	timestamp := uc.now().UTC().Format(time.RFC3339)
	name := domain.NormalizeName(input.Name)
	message := domain.BuildGreeting(prefix, name)
	record := domain.GreetingRecord{
		RequestID: requestID,
		Name:      name,
		Message:   message,
		CreatedAt: timestamp,
		Source:    source,
	}

	if err := uc.recorder.Record(ctx, record); err != nil {
		return GreetingOutput{}, apperrors.Wrap(apperrors.ErrIntegration, "failed to record greeting request", err)
	}

	uc.log.Info("greeting served", "request_id", requestID, "name", name, "source", source)

	return GreetingOutput{
		Message:   message,
		RequestID: requestID,
		Source:    source,
		Timestamp: timestamp,
	}, nil
}
