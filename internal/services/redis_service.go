package services

import (
	"context"
	"time"
	"trieu_mock_project_go/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService() *RedisService {
	return &RedisService{
		client: config.RedisClient,
	}
}

func (s *RedisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.client.Set(ctx, key, value, expiration).Err()
}

func (s *RedisService) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *RedisService) Del(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

func (s *RedisService) Exists(ctx context.Context, key string) (bool, error) {
	n, err := s.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (s *RedisService) Publish(ctx context.Context, channel string, message interface{}) error {
	return s.client.Publish(ctx, channel, message).Err()
}

func (s *RedisService) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return s.client.Subscribe(ctx, channel)
}
