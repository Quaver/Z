package db

import (
	"example.com/Quaver/Z/common"
	"fmt"
)

type UserStats struct {
	Mode                     common.Mode
	Rank                     int
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
func GetUserStats(userId int, mode common.Mode) (*UserStats, error) {
	modeStr, err := common.GetModeString(mode)

	if err != nil {
		return nil, err
	}

	table := fmt.Sprintf("user_stats_%v", modeStr)
	query := fmt.Sprintf("SELECT * FROM %v WHERE user_id = ?", table)
	stats := UserStats{Mode: mode, Rank: -1}

	err = SQL.Get(&stats, query, userId)

	if err != nil {
		return nil, err
	}

	// TODO: Get user rank from the database.

	return &stats, nil
}
