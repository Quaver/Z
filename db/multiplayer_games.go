package db

import (
	"time"
)

// InsertMultiplayerGame Inserts a multiplayer game into the database. Returns the id of the game
func InsertMultiplayerGame(name string, uniqueGameId string) (int, error) {
	query := "INSERT INTO multiplayer_games (unique_id, name, type, time_created) VALUES (?, ?, ?, ?)"
	result, err := SQL.Exec(query, uniqueGameId, name, 0, time.Now().UnixMilli())

	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	return int(id), nil
}
