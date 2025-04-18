package multiplayer

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/objects"
	"example.com/Quaver/Z/scoring"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
	"fmt"
	"log"
	"strconv"
)

// ClearRedisGames Clears all cached multiplayer games in Redis (usually done once at server start)
func ClearRedisGames() error {
	err := db.ClearRedisKeysWithPattern("quaver:server:multiplayer:*")

	if err != nil {
		return err
	}

	return err
}

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
		"ep", strconv.Itoa(utils.BoolToInt(game.Data.EnablePreview)),
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

// Caches a player in Redis
func (game *Game) cachePlayer(id int) {
	user := sessions.GetUserById(id)

	if user == nil {
		return
	}

	wins, err := utils.Find(game.Data.PlayerWins, func(x *objects.MultiplayerGamePlayerWins) bool { return x.Id == id })

	if err != nil {
		wins = &objects.MultiplayerGamePlayerWins{Wins: 0}
	}

	mods, err := utils.Find(game.Data.PlayerModifiers, func(x *objects.MultiplayerGamePlayerMods) bool { return x.Id == id })

	if err != nil {
		mods = &objects.MultiplayerGamePlayerMods{Modifiers: 0}
	}

	player := []string{
		"id", strconv.Itoa(user.Info.Id),
		"u", user.Info.Username,
		"sid", user.Info.SteamId,
		"a", user.Info.AvatarUrl.String,
		"c", user.Info.Country,
		"w", strconv.Itoa(wins.Wins),
		"m", strconv.Itoa(int(mods.Modifiers)),
		"r", strconv.Itoa(utils.BoolToInt(utils.Includes(game.Data.PlayersReady, id))),
		"hm", strconv.Itoa(utils.BoolToInt(!utils.Includes(game.Data.PlayersWithoutMap, id))),
		// "t", strconv.Itoa(0) - Team
	}

	_, err = db.Redis.HSet(db.RedisCtx, game.getPlayerRedisKey(id), player).Result()

	if err != nil {
		log.Printf("Failed to cache multiplayer player in redis - %v\n", err)
		return
	}
}

// Deletes a cached player in redis
func (game *Game) deleteCachedPlayer(userId int) {
	_, err := db.Redis.Del(db.RedisCtx, game.getPlayerRedisKey(userId), game.getPlayerScoreRedisKey(userId)).Result()

	if err != nil {
		log.Printf("Failed to remove multiplayer player in redis - %v\n", err)
		return
	}
}

// Returns the redis key for a player's score in redis.
func (game *Game) getPlayerScoreRedisKey(userId int) string {
	return fmt.Sprintf("quaver:server:multiplayer:%v:%v", game.Data.GameId, userId)
}

// Caches a player's score in redis.
func (game *Game) cachePlayerScore(userId int, processor *scoring.ScoreProcessor) {
	player := []string{
		"m", strconv.FormatInt(int64(processor.Modifiers), 10),
		"pr", strconv.FormatFloat(processor.PerformanceRating, 'f', -1, 64),
		"a", strconv.FormatFloat(processor.Accuracy, 'f', -1, 64),
		"mc", strconv.Itoa(processor.MaxCombo),
		"ma", strconv.Itoa(processor.Judgements[common.JudgementMarv]),
		"pf", strconv.Itoa(processor.Judgements[common.JudgementPerf]),
		"gr", strconv.Itoa(processor.Judgements[common.JudgementGreat]),
		"gd", strconv.Itoa(processor.Judgements[common.JudgementGood]),
		"ok", strconv.Itoa(processor.Judgements[common.JudgementOkay]),
		"ms", strconv.Itoa(processor.Judgements[common.JudgementMiss]),
		"cm", strconv.Itoa(processor.Combo),
		// "t", "0", - Team
		// "sc", "0", - Score
		// "hl", strconv.Itoa(100), - Health
		// "fc", "0" - Boolean for full combo
		// "lv" - Lives
		// "hf" - Has Failed
		// "rh" - Is Regenerating Health
		// "br" - Battle Royale Rank
	}

	_, err := db.Redis.HSet(db.RedisCtx, game.getPlayerScoreRedisKey(userId), player).Result()

	if err != nil {
		log.Printf("Failed to cache multiplayer player score in redis - %v\n", err)
		return
	}

}
