package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/0xBoji/web3-edu-core/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

// Setup initializes the Redis client
func Setup() {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisSetting.Host, config.RedisSetting.Port),
		Password: config.RedisSetting.Password,
		DB:       config.RedisSetting.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection established")
}

// GetClient returns the Redis client
func GetClient() *redis.Client {
	return Client
}
