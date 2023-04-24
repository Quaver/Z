package db

import (
	"example.com/Quaver/Z/common"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
)

type UserStats struct {
	Mode                     common.Mode
	GlobalRank               int
	CountryRank              int
	UserId                   int     `db:"user_id"`
	TotalScore               int64   `db:"total_score"`
	RankedScore              int64   `db:"ranked_score"`
	OverallAccuracy          float64 `db:"overall_accuracy"`
	OverallPerformanceRating float64 `db:"overall_performance_rating"`
	PlayCount                int     `db:"play_count"`
	FailCount                int     `db:"fail_count"`
	MaxCombo                 int     `db:"max_combo"`
	ReplaysWatched           int     `db:"replays_watched"`
	TotalMarv                int     `db:"total_marv"`
	TotalPerf                int     `db:"total_perf"`
	TotalGreat               int     `db:"total_great"`
	TotalGood                int     `db:"total_good"`
	TotalOkay                int     `db:"total_okay"`
	TotalMiss                int     `db:"total_miss"`
	TotalPauses              int     `db:"total_pauses"`
	MultiplayerWins          int     `db:"multiplayer_wins"`
	MultiplayerLosses        int     `db:"multiplayer_losses"`
	MultiplayerTies          int     `db:"multiplayer_ties"`
	CountGradeX              int     `db:"count_grade_x"`
	CountGradeSS             int     `db:"count_grade_ss"`
	CountGradeS              int     `db:"count_grade_s"`
	CountGradeA              int     `db:"count_grade_a"`
	CountGradeB              int     `db:"count_grade_b"`
	CountGradeC              int     `db:"count_grade_c"`
	CountGradeD              int     `db:"count_grade_d"`
}

// GetUserStats Fetches the user stats for a given game mode from the database.
func GetUserStats(userId int, country string, mode common.Mode) (*UserStats, error) {
	modeStr, err := common.GetModeString(mode)

	if err != nil {
		return nil, err
	}

	table := fmt.Sprintf("user_stats_%v", modeStr)
	query := fmt.Sprintf("SELECT * FROM %v WHERE user_id = ?", table)
	stats := UserStats{Mode: mode, GlobalRank: -1, CountryRank: -1}

	err = SQL.Get(&stats, query, userId)

	if err != nil {
		return nil, err
	}

	stats.GlobalRank, err = GetUserGlobalRank(userId, mode)

	if err != nil {
		return nil, err
	}

	stats.CountryRank, err = GetUserCountryRank(userId, country, mode)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetUserGlobalRank Retrieves a user's global rank from Redis. If it doesn't exist, it'll return -1.
func GetUserGlobalRank(userId int, mode common.Mode) (int, error) {
	key := fmt.Sprintf("quaver:leaderboard:%v", int(mode))

	rank, err := getRankFromRedis(key, userId)

	if err != nil {
		return -1, err
	}

	return rank, nil
}

// GetUserCountryRank Retrieves a user's country rank from Redis. If it doesn't exist, it'll return -1.
func GetUserCountryRank(userId int, country string, mode common.Mode) (int, error) {
	key := fmt.Sprintf("quaver:country_leaderboard:%v:%v", strings.ToLower(country), int(mode))

	rank, err := getRankFromRedis(key, userId)

	if err != nil {
		return -1, err
	}

	return rank, nil
}

// Gets a rank value in Redis from the database.
func getRankFromRedis(key string, userId int) (int, error) {
	result, err := Redis.ZRevRank(RedisCtx, key, strconv.Itoa(userId)).Result()

	if err != nil {
		// Rank does not exist in the database
		if err == redis.Nil {
			return -1, nil
		}

		return -1, err
	}

	// Adding 1 here to get the actual rank because Redis is zero indexed.
	return int(result) + 1, nil
}
