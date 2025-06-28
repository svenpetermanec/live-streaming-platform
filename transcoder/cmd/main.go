package main

import (
	"context"

	"github.com/redis/go-redis/v9"

	"transcoder/internal/guard"
	"transcoder/internal/logging"
	"transcoder/transcoder/cmd/config"
	"transcoder/transcoder/internal"
)

func main() {
	config.Load()

	logger := logging.NewLogger(config.Cfg.LogLevel)

	guard.CapturePanic(logger)

	redisClient := redis.NewClient(
		&redis.Options{
			Addr: config.Cfg.Redis.Address,
			DB:   config.Cfg.Redis.Database,
		},
	)

	srtServer := internal.NewSrtServer(
		internal.NewSrtManager(
			internal.NewRepository(redisClient),
			logger,
		),
		logger,
		config.Cfg.SRT,
	)

	srtServer.Start(context.Background())
}
