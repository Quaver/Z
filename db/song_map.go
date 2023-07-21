package db

import (
	"database/sql"
	"example.com/Quaver/Z/common"
)

type SongMap struct {
	Id               int            `db:"id"`
	MapsetId         int            `db:"mapset_id"`
	Md5              sql.NullString `db:"md5"`
	AlternativeMd5   sql.NullString `db:"alternative_md5"`
	GameMode         common.Mode    `db:"game_mode"`
	Artist           sql.NullString `db:"artist"`
	Title            sql.NullString `db:"title"`
	DifficultyName   sql.NullString `db:"difficulty_name"`
	DifficultyRating float64        `db:"difficulty_rating"`
}

// GetSongMapById Retrieves a map by its id.
func GetSongMapById(id int) (*SongMap, error) {
	query := "SELECT id, mapset_id, md5, alternative_md5, game_mode, artist, title, difficulty_name, difficulty_rating " +
		"FROM maps WHERE id = ? LIMIT 1"

	var songMap SongMap

	err := SQL.Get(&songMap, query, id)

	if err != nil {
		return nil, err
	}

	return &songMap, nil
}
