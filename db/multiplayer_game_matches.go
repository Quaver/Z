package db

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/objects"
)

type MultiplayerMatch struct {
	Id              int
	GameId          int                            `db:"game_id"`
	TimePlayed      int64                          `db:"time_played"`
	MapMd5          string                         `db:"map_md5"`
	MapName         string                         `db:"map"`
	HostId          int                            `db:"host_id"`
	Ruleset         objects.MultiplayerGameRuleset `db:"ruleset"`
	GameMode        common.Mode                    `db:"game_mode"`
	GlobalModifiers common.Mods                    `db:"global_modifiers"`
	FreeMod         objects.MultiplayerGameFreeMod `db:"free_mod_type"`
	HealthType      objects.MultiplayerGameHealth  `db:"health_type"`
	Lives           int                            `db:"lives"`
	Aborted         bool                           `db:"aborted"`
}

// InsertIntoDatabase Inserts a multiplayer match into the database and returns the insert id of it.
func (match *MultiplayerMatch) InsertIntoDatabase() error {
	query := "INSERT INTO multiplayer_game_matches" +
		"(game_id, time_played, map_md5, map, host_id, ruleset, game_mode, global_modifiers, free_mod_type, health_Type, lives, aborted) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	result, err := SQL.Exec(query, match.GameId, match.TimePlayed, match.MapMd5, match.MapName, match.HostId, match.Ruleset, match.GameMode,
		match.GlobalModifiers, match.FreeMod, match.HealthType, match.Lives, match.Aborted)

	if err != nil {
		match.Id = -1
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		match.Id = -1
		return err
	}

	match.Id = int(id)
	return nil
}
