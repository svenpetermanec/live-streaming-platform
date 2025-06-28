package internal

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) SetStreamName(ctx context.Context, streamId string, streamName string) error {
	return r.client.Set(ctx, streamId, streamName, 0).Err()
}
