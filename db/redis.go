package db

import (
	"context"
	"example.com/Quaver/Z/config"
	"github.com/go-redis/redis/v8"
	"log"
)

var (
	Redis                            *redis.Client
	RedisCtx                         = context.Background()
	redisChannelHandlers             = map[string][]func(message *redis.Message){}
	RedisChannelSongRequests         = "quaver:song_requests"
	RedisChannelTwitchConnection     = "quaver:twitch_connection"
	RedisChannelMultiplayerMapShares = "quaver:multiplayer_map_shares"
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
	log.Println("Successfully connected to redis")

	if result.Err() != nil {
		log.Fatalln(result.Err())
	}

	sub := Redis.Subscribe(RedisCtx, RedisChannelSongRequests, RedisChannelTwitchConnection, RedisChannelMultiplayerMapShares)

	go func() {
		for {
			msg, err := sub.ReceiveMessage(RedisCtx)

			if err != nil {
				log.Printf("Error receiving redis message - %v", err)
			}

			// Go through and call all the handler functions for this particular channel.
			if handlers, ok := redisChannelHandlers[msg.Channel]; ok {
				for _, handler := range handlers {
					handler(msg)
				}
			}
		}
	}()
}

// AddRedisSubscriberHandler Adds a handler to a given channel
func AddRedisSubscriberHandler(channel string, f func(message *redis.Message)) {
	if _, ok := redisChannelHandlers[channel]; !ok {
		redisChannelHandlers[channel] = []func(message *redis.Message){}
	}

	redisChannelHandlers[channel] = append(redisChannelHandlers[channel], f)
}

// ClearRedisKeysWithPattern Clears a given pattern of redis keys from the database
func ClearRedisKeysWithPattern(pattern string) error {
	keys, err := Redis.Keys(RedisCtx, pattern).Result()

	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	_, err = Redis.Del(RedisCtx, keys...).Result()

	if err != nil {
		return err
	}

	return nil
}
