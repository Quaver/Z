package multiplayer

import (
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/utils"
	"fmt"
	"log"
	"strconv"
)

// Returns the redis key for the match settings
func (game *Game) getMatchSettingsRedisKey() string {
	return fmt.Sprintf("quaver:server:multiplayer:%v", game.Data.Id)
}

// Caches the current match settings in redis
func (game *Game) cacheMatchSettings() {
	settings := []string{
		"n", game.Data.Name,
		"pw", strconv.Itoa(utils.BoolToInt(game.Data.HasPassword)),
		"mp", strconv.Itoa(game.Data.MaxPlayers),
		"md5", game.Data.MapMD5,
		"mid", strconv.Itoa(game.Data.MapId),
		"msid", strconv.Itoa(game.Data.MapsetId),
		"map", game.Data.MapName,
		"host", strconv.Itoa(game.Data.HostId),
		"r", strconv.Itoa(int(game.Data.Ruleset)),
		"hr", strconv.Itoa(utils.BoolToInt(game.Data.IsHostRotation)),
		"gm", strconv.Itoa(int(game.Data.MapGameMode)),
		"d", strconv.FormatFloat(game.Data.MapDifficultyRating, 'f', -1, 64),
		"inp", strconv.Itoa(utils.BoolToInt(game.Data.InProgress)),
		"m", strconv.FormatInt(int64(game.Data.GlobalModifiers), 10),
		"fm", strconv.Itoa(int(game.Data.FreeModType)),
		"trn", strconv.Itoa(utils.BoolToInt(game.Data.IsTournamentMode)),
		// "t", strconv.Itoa(0), -  Game Type
		// "h", strconv.Itoa(0), - Health Type
		// "lv", strconv.Itoa(3) - Life Count
		// "rtw", strconv.Itoa(game.Data.TeamRedWins), - Red Team Wins
		// "btw", strconv.Itoa(game.DAta.TeamBlueWins), - Blue Team Wins
	}

	_, err := db.Redis.HSet(db.RedisCtx, game.getMatchSettingsRedisKey(), settings).Result()

	if err != nil {
		log.Printf("Failed to cache match settings in redis - %v\n", err)
		return
	}
}

// Deletes the cached match settings in redis
func (game *Game) deleteCachedMatchSettings() {
	_, err := db.Redis.Del(db.RedisCtx, game.getMatchSettingsRedisKey()).Result()

	if err != nil {
		log.Printf("Failed to remove match settings in redis - %v\n", err)
		return
	}
}

// Returns the redis key for an individual user in the game
func (game *Game) getPlayerRedisKey(id int) string {
	return fmt.Sprintf("quaver:server:multiplayer:%v:player:%v", game.Data.Id, id)
}

// Returns the redis key for a player's score in the game
func (game *Game) getPlayerScoreRedisKey(id int) string {
	return fmt.Sprintf("quaver:server:multiplayer:%v:%v", game.Data.Id, id)
}
