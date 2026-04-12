// Package usecase provides business logic use cases.
package usecase

import (
	"github.com/PlatformStackPulse/go-template/internal/domain"
	"github.com/PlatformStackPulse/go-template/internal/logger"
)

// GreetingUseCase handles greeting logic
type GreetingUseCase struct {
	log *logger.Logger
}

// NewGreetingUseCase creates a new greeting use case
func NewGreetingUseCase(log *logger.Logger) *GreetingUseCase {
	return &GreetingUseCase{log: log}
}

// Execute executes the greeting use case
func (uc *GreetingUseCase) Execute(name string) (string, error) {
	greeter := domain.NewGreeter(name)
	message := greeter.Greet()
	uc.log.Debug("Greeting generated", "name", name, "message", message)
	return message, nil
}
