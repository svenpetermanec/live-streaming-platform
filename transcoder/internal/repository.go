package internal

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) GetStreamName(ctx context.Context, streamId string) (string, error) {
	streamName, err := r.client.Get(ctx, streamId).Result()
	if err != nil {
		return "", err
	}

	return streamName, nil
}
