package main

import (
	"context"

	"gitlab.com/mildd/flow-gateway/internal/app"
	"gitlab.com/mildd/flow-gateway/internal/configs"
	"gitlab.com/mildd/flow-gateway/pkg/logging"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logging.Info(ctx, "config initializing")
	cfg, err := configs.GetConfig()
	if err != nil {
		logging.Error(ctx, err)

	}

	ctx = logging.ContextWithLogger(ctx, logging.NewLogger())

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		logging.Fatal(ctx, err)
	}

	logging.Info(ctx, "Running Application")
	a.Run(ctx)
}
