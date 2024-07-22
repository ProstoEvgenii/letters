package rediscache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"mailsender/config"
	"time"
)

type Client struct {
	Client *redis.Client
}

// NewRedisClient returns new redis client
func NewRedisClient(cfg *config.Config) (Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	redisHost := cfg.Redis.Address

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := client.Ping(ctx).Result()

	return Client{client}, err
}
