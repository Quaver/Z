package db

import "database/sql"

type Clan struct {
	Id           int    `db:"id"`
	OwnerId      int    `db:"owner_id"`
	Name         string `db:"name"`
	Tag          string `db:"tag"`
	Customizable bool   `db:"customizable"`
}

// GetAllClans Retrieves all of the clans in the db
func GetAllClans() ([]*Clan, error) {
	result := make([]*Clan, 0)

	err := SQL.Select(&result, "SELECT id, owner_id, name, tag, customizable FROM clans")

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return result, nil
}
