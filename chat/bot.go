package chat

import (
	"database/sql"
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/sessions"
	"fmt"
	"strings"
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

// handleBotCommands Handles bot commands for a given message (if any)
func handleBotCommands(user *sessions.User, channel *Channel, message string) {
	if message == "" || message[0] != '!' {
		return
	}

	args := strings.Split(message, " ")

	if strings.Contains(message, "!mp") {
		handleMultiplayerCommands(user, channel, args)
	}
}

// Handles commands made for multiplayer
func handleMultiplayerCommands(user *sessions.User, channel *Channel, args []string) {
	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	fmt.Println("USER EXECUTED COMMAND ", args)
}
