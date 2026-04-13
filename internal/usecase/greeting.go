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
	// Keys sourced from SSM app-config and environment overrides.
	sampleGreetingPrefixKey    = "sample.greeting.prefix"
	sampleGreetingPrefixEnvKey = "SAMPLE_GREETING_PREFIX"

	// Keys for response source normalization.
	apiSourceLabelEnvKey       = "API_SOURCE_LABEL"
	defaultGreetingSourceLabel = "api"
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

	prefix, err := uc.resolveGreetingPrefix(ctx)
	if err != nil {
		return GreetingOutput{}, err
	}

	requestID := normalizeRequestID(input.RequestID)
	source := uc.resolveSourceLabel(ctx, input.Source)

	timestamp := uc.now().UTC().Format(time.RFC3339)
	name := domain.NormalizeName(input.Name)
	message := domain.BuildGreeting(prefix, name)
	record := newGreetingRecord(requestID, name, message, timestamp, source)

	// Persist request metadata so teams can inspect sample traffic in backing stores.
	if err := uc.recorder.Record(ctx, record); err != nil {
		return GreetingOutput{}, apperrors.Wrap(apperrors.ErrIntegration, "failed to record greeting request", err)
	}

	uc.log.Info("greeting served", "request_id", requestID, "name", name, "source", source)

	return newGreetingOutput(message, requestID, source, timestamp), nil
}

func (uc *GreetingUseCase) resolveGreetingPrefix(ctx context.Context) (string, error) {
	// Environment overrides are checked first for fast local/test iteration.
	if uc.environment != nil {
		if prefix := uc.environment.Lookup(ctx, sampleGreetingPrefixEnvKey); prefix != "" {
			return prefix, nil
		}
	}

	// Default source of truth is the shared app-config document in SSM.
	prefix, err := uc.configProvider.StringValue(ctx, sampleGreetingPrefixKey)
	if err != nil {
		return "", apperrors.Wrap(apperrors.ErrIntegration, "failed to load greeting prefix", err)
	}

	return prefix, nil
}

func (uc *GreetingUseCase) resolveSourceLabel(ctx context.Context, inputSource string) string {
	source := inputSource
	if source == "" {
		source = defaultGreetingSourceLabel
	}

	// API source label can be normalized by environment to simplify observability.
	if uc.environment != nil {
		if sourceLabel := uc.environment.Lookup(ctx, apiSourceLabelEnvKey); sourceLabel != "" {
			return sourceLabel
		}
	}

	return source
}

func normalizeRequestID(inputRequestID string) string {
	if inputRequestID == "" {
		return "unknown"
	}

	return inputRequestID
}

func newGreetingRecord(requestID, name, message, timestamp, source string) domain.GreetingRecord {
	return domain.GreetingRecord{
		RequestID: requestID,
		Name:      name,
		Message:   message,
		CreatedAt: timestamp,
		Source:    source,
	}
}

func newGreetingOutput(message, requestID, source, timestamp string) GreetingOutput {
	return GreetingOutput{
		Message:   message,
		RequestID: requestID,
		Source:    source,
		Timestamp: timestamp,
	}
}
