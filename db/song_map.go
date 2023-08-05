package db

import (
	"database/sql"
	"errors"
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

// GetRandomSongMap Retrieves a random map from the database with min/max difficulty rating filter
func GetRandomSongMap(minDiff float32, maxDiff float32, gameModes []common.Mode) (*SongMap, error) {
	mode := ""

	switch len(gameModes) {
	case 0:
		return nil, errors.New("no game modes provided")
	case 1:
		mode = "AND game_mode = ? "
	}

	query := "SELECT id, mapset_id, md5, alternative_md5, game_mode, artist, title, difficulty_name, difficulty_rating " +
		"FROM maps " +
		"WHERE difficulty_rating > ? AND difficulty_rating < ? AND ranked_status = 2 " + mode +
		"ORDER BY RAND() " +
		"LIMIT 1"

	var songMap SongMap
	var err error

	if mode != "" {
		err = SQL.Get(&songMap, query, minDiff, maxDiff, gameModes[0])
	} else {
		err = SQL.Get(&songMap, query, minDiff, maxDiff)
	}

	if err != nil {
		return nil, err
	}

	return &songMap, nil
}
