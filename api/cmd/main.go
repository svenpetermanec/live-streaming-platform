package main

import (
	"context"

	"transcoder/api/cmd/config"
	"transcoder/api/internal"
	"transcoder/internal/guard"
	"transcoder/internal/logging"
)

func main() {
	config.Load()

	logger := logging.NewLogger(config.Cfg.LogLevel)

	guard.CapturePanic(logger)

	httpServer := internal.NewHttpServer(
		internal.NewController(config.Cfg, logger),
		config.Cfg.HTTPServer.Port,
		logger,
	)

	httpServer.Start(context.Background())
}
