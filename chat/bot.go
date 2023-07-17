package chat

import (
	"database/sql"
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/sessions"
)

var (
	Bot = sessions.NewUser(nil, &db.User{
		Id:             2,
		SteamId:        "1",
		Username:       "Quaver",
		Allowed:        true,
		Privileges:     common.PrivilegeNormal,
		UserGroups:     common.UserGroupBot | common.UserGroupNormal,
		MuteEndTime:    0,
		Country:        "US",
		AvatarUrl:      "",
		TwitchUsername: sql.NullString{},
	})
)
