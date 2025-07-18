package services

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type cacheService struct {
	client *redis.Client
}

func NewCacheService(client *redis.Client) CacheService {
	return &cacheService{
		client: client,
	}
}

func (c *cacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *cacheService) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (c *cacheService) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
