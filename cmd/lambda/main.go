package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/PlatformStackPulse/go-lambda-template/internal/app"
)

func main() {
	application, err := app.New(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to bootstrap lambda application: %v\n", err)
		os.Exit(1)
	}

	lambda.Start(application.Handler.Handle)
}
