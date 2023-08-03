package db

import (
	"database/sql"
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/config"
	"fmt"
	"github.com/Philipp15b/go-steamapi"
	"strconv"
	"time"
)

type User struct {
	Id             int               `db:"id"`
	SteamId        string            `db:"steam_id"`
	Username       string            `db:"username"`
	Allowed        bool              `db:"allowed"`
	Privileges     common.Privileges `db:"privileges"`
	UserGroups     common.UserGroups `db:"usergroups"`
	MuteEndTime    int64             `db:"mute_endtime"`
	Country        string            `db:"country"`
	AvatarUrl      sql.NullString    `db:"avatar_url"`
	TwitchUsername sql.NullString    `db:"twitch_username"`
}

// GetProfileUrl Returns the full profile url for the user
func (u *User) GetProfileUrl() string {
	return fmt.Sprintf("https://quavergame.com/user/%v", u.Id)
}

// GetUserBySteamId Retrieves a user from the database by their Steam id
func GetUserBySteamId(steamId string) (*User, error) {
	query := "SELECT id, steam_id, username, allowed, privileges, usergroups, mute_endtime, country, avatar_url, twitch_username FROM users WHERE steam_id = ? LIMIT 1"

	var user User
	err := SQL.Get(&user, query, steamId)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername Rerieves a user from the database by their username
func GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, steam_id, username, allowed, privileges, usergroups, mute_endtime, country, avatar_url, twitch_username FROM users WHERE username = ? LIMIT 1"

	var user User
	err := SQL.Get(&user, query, username)

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

// UpdateUserSteamAvatar Updates the Steam avatar_url of a given user.
// Returns a link to the avatar
func UpdateUserSteamAvatar(steamId string) (string, error) {
	parsedId, err := strconv.ParseInt(steamId, 10, 64)

	if err != nil {
		panic(err)
	}

	summaries, err := steamapi.GetPlayerSummaries([]uint64{uint64(parsedId)}, config.Instance.Steam.PublisherKey)

	if err != nil {
		return "", err
	}

	if len(summaries) == 0 {
		return "", fmt.Errorf("steam player summaries returned 0 users")
	}

	avatar := summaries[0].LargeAvatarURL

	_, err = SQL.Exec("UPDATE users SET avatar_url = ? WHERE steam_id = ?", avatar, steamId)

	if err != nil {
		return "", err
	}

	return avatar, nil
}

// MuteUser Mutes a user for a given duration
func MuteUser(id int, endTime int64) error {
	_, err := SQL.Exec("UPDATE users SET mute_endtime = ? WHERE id = ?", endTime, id)

	if err != nil {
		return err
	}

	return nil
}

// UnlinkUserTwitch Unlinks the twitch account of a given user
func UnlinkUserTwitch(id int) error {
	_, err := SQL.Exec("UPDATE users SET twitch_username = NULL WHERE id = ?", id)

	if err != nil {
		return err
	}

	return nil
}
