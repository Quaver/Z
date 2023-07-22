package handlers

import (
	"encoding/json"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"github.com/go-redis/redis/v8"
	"log"
)

type RedisTwitchSongRequest struct {
	UserId  int `json:"user_id"`
	Request struct {
		TwitchUsername   string                  `json:"twitch_username"`
		Game             packets.SongRequestGame `json:"game"`
		MapId            int                     `json:"map_id"`
		MapsetId         int                     `json:"mapset_id"`
		MapMd5           string                  `json:"map_md5"`
		Artist           string                  `json:"artist"`
		Title            string                  `json:"title"`
		DifficultyName   string                  `json:"difficulty_name"`
		Creator          string                  `json:"creator"`
		DifficultyRating float64                 `json:"difficulty_rating"`
	} `json:"request"`
}

func AddRedisHandlers() {
	db.AddRedisSubscriberHandler(db.RedisChannelSongRequests, HandleTwitchSongRequest)
	db.AddRedisSubscriberHandler(db.RedisChannelTwitchConnection, HandleTwitchConnection)
	db.AddRedisSubscriberHandler(db.RedisChannelMultiplayerMapShares, HandleMultiplayerMapShares)
}

func HandleTwitchSongRequest(msg *redis.Message) {
	var parsed RedisTwitchSongRequest

	err := json.Unmarshal([]byte(msg.Payload), &parsed)

	if err != nil {
		log.Printf("Failed to parse twitch song request - %v - %v", msg.Payload, err)
	}

	user := sessions.GetUserById(parsed.UserId)

	if user == nil {
		return
	}

	sessions.SendPacketToUser(packets.NewServerSongRequest(packets.SongRequest{
		TwitchUsername:   parsed.Request.TwitchUsername,
		UserId:           -1,
		Game:             parsed.Request.Game,
		MapId:            parsed.Request.MapId,
		MapsetId:         parsed.Request.MapsetId,
		MapMd5:           parsed.Request.MapMd5,
		Artist:           parsed.Request.Artist,
		Title:            parsed.Request.Title,
		DifficultyName:   parsed.Request.DifficultyName,
		Creator:          parsed.Request.Creator,
		DifficultyRating: parsed.Request.DifficultyRating,
	}), user)
}

func HandleTwitchConnection(msg *redis.Message) {

}

func HandleMultiplayerMapShares(msg *redis.Message) {

}
