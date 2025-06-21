package main

import (
	"context"

	"transcoder/internal/guard"
	"transcoder/internal/logging"
	"transcoder/transcoder/cmd/config"
	"transcoder/transcoder/internal"
)

func main() {
	config.Load()

	logger := logging.NewLogger(config.Cfg.LogLevel)

	guard.CapturePanic(logger)

	srtServer := internal.NewSrtServer(
		internal.NewSrtManager(logger),
		logger,
		config.Cfg.SRT,
	)

	srtServer.Start(context.Background())
}
