package handlers

import (
	"example.com/Quaver/Z/db"
	"github.com/go-redis/redis/v8"
)

func AddRedisHandlers() {
	db.AddRedisSubscriberHandler(db.RedisChannelSongRequests, HandleTwitchSongRequest)
	db.AddRedisSubscriberHandler(db.RedisChannelTwitchConnection, HandleTwitchConnection)
	db.AddRedisSubscriberHandler(db.RedisChannelMultiplayerMapShares, HandleMultiplayerMapShares)
}

func HandleTwitchSongRequest(msg *redis.Message) {

}

func HandleTwitchConnection(msg *redis.Message) {

}

func HandleMultiplayerMapShares(msg *redis.Message) {

}
