package handlers

import (
	"encoding/json"
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

func AddRedisHandlers() {
	db.AddRedisSubscriberHandler(db.RedisChannelSongRequests, HandleTwitchSongRequest)
	db.AddRedisSubscriberHandler(db.RedisChannelTwitchConnection, HandleTwitchConnection)
	db.AddRedisSubscriberHandler(db.RedisChannelMultiplayerMapShares, HandleMultiplayerMapShares)
	db.AddRedisSubscriberHandler(db.RedisChannelFirstPlaceScores, HandleFirstPlaceScores)
}

func HandleTwitchSongRequest(msg *redis.Message) {
	type redisTwitchSongRequest struct {
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

	var parsed redisTwitchSongRequest

	err := json.Unmarshal([]byte(msg.Payload), &parsed)

	if err != nil {
		log.Printf("Failed to parse twitch song request - %v - %v\n", msg.Payload, err)
		return
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
	type redisTwitchConnection struct {
		UserId int `json:"user_id"`
	}

	var parsed redisTwitchConnection

	err := json.Unmarshal([]byte(msg.Payload), &parsed)

	if err != nil {
		log.Printf("Failed to parse redis twitch connection - %v - %v\n", msg.Payload, err)
		return
	}

	user := sessions.GetUserById(parsed.UserId)

	if user == nil {
		return
	}

	newUser, err := db.GetUserBySteamId(user.Info.SteamId)

	if err != nil {
		log.Printf("Failed to retrieve user from DB while handling redis twitch connection - %v\n", err)
		return
	}

	user.Info.TwitchUsername = newUser.TwitchUsername
	sessions.SendPacketToUser(packets.NewServerTwitchConnection(user.Info.TwitchUsername.String), user)
}

func HandleMultiplayerMapShares(msg *redis.Message) {
	type redisMultiplayerMapShare struct {
		UploaderId int `json:"uploader_id"`
		GameId     int `json:"game_id"`
	}

	var parsed redisMultiplayerMapShare

	err := json.Unmarshal([]byte(msg.Payload), &parsed)

	if err != nil {
		log.Printf("Failed to parse redis multiplayer map share - %v - %v\n", msg.Payload, err)
		return
	}

	user := sessions.GetUserById(parsed.UploaderId)

	if user == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.RunLocked(func() {
		game.SetDonatorMapsetShared(true, true)
	})
}

func HandleFirstPlaceScores(msg *redis.Message) {
	type redisFirstPlaceScore struct {
		User struct {
			Username string `json:"username"`
		} `json:"user"`
		Map struct {
			Artist         string `json:"artist"`
			Title          string `json:"title"`
			DifficultyName string `json:"difficulty_name"`
		} `json:"map"`
	}

	var parsed redisFirstPlaceScore

	err := json.Unmarshal([]byte(msg.Payload), &parsed)

	if err != nil {
		log.Printf("Failed to parse first place score - %v - %v\n", msg.Payload, err)
		return
	}

	chat.SendMessage(chat.Bot, "#first-places", fmt.Sprintf("%v has just achieved first place on %v - %v [%v]",
		parsed.User.Username, parsed.Map.Artist, parsed.Map.Title, parsed.Map.DifficultyName))
}
