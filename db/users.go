package db

import "time"

type User struct {
	Id          int    `db:"id"`
	SteamId     string `db:"steam_id"`
	Username    string `db:"username"`
	Allowed     bool   `db:"allowed"`
	Privileges  int64  `db:"privileges"` // TODO: USE ENUM
	UserGroups  int64  `db:"usergroups"` // TODO: USE ENUM
	MuteEndTime int64  `db:"mute_endtime"`
	Country     string `db:"country"`
	AvatarUrl   string `db:"avatar_url"`
}

// GetUserBySteamId Retrieves a user from the database by their Steam id
func GetUserBySteamId(id string) (*User, error) {
	var user User
	err := SQL.Get(&user, "SELECT id, steam_id, username, allowed, privileges, usergroups, mute_endtime, country, avatar_url FROM users WHERE steam_id = ? LIMIT 1", id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserLatestActivity Updates the latest_activity of a user to the current time
func UpdateUserLatestActivity(id int) error {
	_, err := SQL.Exec("UPDATE users SET latest_activity = ? WHERE id = ?", time.Now().UnixMilli(), id)

	if err != nil {
		return err
	}

	return nil
}
