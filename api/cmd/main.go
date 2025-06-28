package main

import (
	"context"

	"github.com/redis/go-redis/v9"

	"transcoder/api/cmd/config"
	"transcoder/api/internal"
	"transcoder/internal/guard"
	"transcoder/internal/logging"
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

	httpServer := internal.NewHttpServer(
		internal.NewController(
			internal.NewRepository(redisClient),
			config.Cfg, logger,
		),
		config.Cfg.HTTPServer.Port,
		logger,
	)

	httpServer.Start(context.Background())
}
