package main

import (
	"context"

	"transcoder/internal/guard"
	"transcoder/internal/logging"
	"transcoder/transcoder/cmd/config"
	internal2 "transcoder/transcoder/internal"
)

func main() {
	config.Load()

	logger := logging.NewLogger()

	guard.CapturePanic(logger)

	srtServer := internal2.NewSrtServer(
		internal2.NewSrtManager(logger),
		logger,
		config.Cfg.SRT,
	)

	srtServer.Start(context.Background())
}
