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
	Id                      int               `db:"id"`
	SteamId                 string            `db:"steam_id"`
	Username                string            `db:"username"`
	Allowed                 bool              `db:"allowed"`
	Privileges              common.Privileges `db:"privileges"`
	UserGroups              common.UserGroups `db:"usergroups"`
	MuteEndTime             int64             `db:"mute_endtime"`
	Country                 string            `db:"country"`
	AvatarUrl               sql.NullString    `db:"avatar_url"`
	TwitchUsername          sql.NullString    `db:"twitch_username"`
	ClanId                  sql.NullInt32     `db:"clan_id"`
	ClanTag                 sql.NullString    `db:"clan_tag"`
	ClanAccentColor         sql.NullString    `db:"clan_accent_color"`
	AccentColorCustomizable bool              `db:"accent_color_customizable"`
	AccentColor             sql.NullString    `db:"accent_color"`
}

// GetAccentColor Returns the user's active accent color, if customizable.
func (u *User) GetAccentColor() string {
	if !u.AccentColorCustomizable || !u.AccentColor.Valid {
		return ""
	}

	return u.AccentColor.String
}

// GetProfileUrl Returns the full profile url for the user
func (u *User) GetProfileUrl() string {
	return fmt.Sprintf("https://quavergame.com/users/%v", u.Id)
}

func getUserBy(column string, value interface{}) (*User, error) {
	query := "SELECT users.id, steam_id, username, allowed, privileges, usergroups, mute_endtime, country, avatar_url, twitch_username, clan_id, clans.tag AS clan_tag, clans.accent_color AS clan_accent_color, users.accent_color_customizable, users.accent_color FROM users LEFT JOIN clans ON users.clan_id = clans.id WHERE %s = ? LIMIT 1"

	var user User
	err := SQL.Get(&user, fmt.Sprintf(query, column), value)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserById Retrieves a user from the database by their id
func GetUserById(id int) (*User, error) {
	return getUserBy("users.id", id)
}

// GetUserBySteamId Retrieves a user from the database by their Steam id
func GetUserBySteamId(steamId string) (*User, error) {
	return getUserBy("users.steam_id", steamId)
}

// GetUserByUsername Rerieves a user from the database by their username
func GetUserByUsername(username string) (*User, error) {
	return getUserBy("users.username", username)
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

	summaries, err := steamapi.GetPlayerSummaries([]uint64{uint64(parsedId)}, config.Instance.Steam.APIKey)

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
