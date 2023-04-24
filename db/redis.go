package db

import (
	"context"
	"example.com/Quaver/Z/config"
	"github.com/go-redis/redis/v8"
	"log"
)

var (
	Redis    *redis.Client
	RedisCtx = context.Background()
)

// InitializeRedis Initializes a Redis client
func InitializeRedis() {
	if Redis != nil {
		return
	}

	Redis = redis.NewClient(&redis.Options{
		Addr:     config.Instance.Redis.Host,
		Password: config.Instance.Redis.Password,
		DB:       config.Instance.Redis.Database,
	})

	result := Redis.Ping(RedisCtx)

	if result.Err() != nil {
		log.Fatalln(result.Err())
	}

	log.Println("Successfully connected to redis")
}
