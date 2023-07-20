package db

import "example.com/Quaver/Z/common"

type MultiplayerMatchScore struct {
	UserId            int         `db:"user_id"`
	MatchId           int         `db:"match_id"`
	Mods              common.Mods `db:"mods"`
	PerformanceRating float64     `db:"performance_rating"`
	Accuracy          float64     `db:"accuracy"`
	MaxCombo          int         `db:"max_combo"`
	CountMarv         int         `db:"count_marv"`
	CountPerf         int         `db:"count_perf"`
	CountGreat        int         `db:"count_great"`
	CountGood         int         `db:"count_good"`
	CountOkay         int         `db:"count_okay"`
	CountMiss         int         `db:"count_miss"`
	Won               int         `db:"won"`
}

// InsertIntoDatabase Inserts the score into the database
func (s *MultiplayerMatchScore) InsertIntoDatabase() error {
	query := "INSERT INTO multiplayer_match_scores " +
		"(user_id, match_id, mods, performance_rating, accuracy, max_combo, count_marv, count_perf, count_great, " +
		"count_good, count_okay, count_miss, won, team, score, has_failed, lives_left, full_combo, battle_royale_place) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := SQL.Exec(query, s.UserId, s.MatchId, s.Mods, s.PerformanceRating, s.Accuracy, s.MaxCombo, s.CountMarv, s.CountPerf, s.CountGreat,
		s.CountGood, s.CountOkay, s.CountMiss, s.Won, 0, 0, 0, 0, 0, 0)

	if err != nil {
		return err
	}

	return nil
}
