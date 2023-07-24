package chat

import (
	"database/sql"
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
	"fmt"
	"log"
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

const (
	botMessageStop = "I'm just trying to help you. Why are you doing this to me?"
)

// Adds chat handlers for bot messages
func addBotChatHandlers() {
	AddPublicMessageHandler(handlePublicChatBotCommands)
	AddPrivateMessageHandler(handlePrivateChatBotCommands)
}

// Handles bot commands for public messages
func handlePublicChatBotCommands(user *sessions.User, channel *Channel, args []string) string {
	return handleBotCommands(user, args)
}

// Handles bot commands for private messages
func handlePrivateChatBotCommands(user *sessions.User, receivingUser *sessions.User, args []string) string {
	// Only handle if the user is talking to the bot directly.
	if receivingUser != Bot {
		return "'"
	}

	return handleBotCommands(user, args)
}

// Handles all bot commands regardless of if they're public or private
func handleBotCommands(user *sessions.User, args []string) string {
	if len(args) == 0 || args[0] == "" || args[0][0] != '!' {
		return ""
	}

	log.Println(strings.ToLower(strings.Split(args[0], "!")[1]))

	switch strings.ToLower(strings.Split(args[0], "!")[1]) {
	case "kick":
		return handleBotCommandKick(user, args)
	case "alertall", "notifyall":
		return handleBotCommandNotifyAll(user, args)
	default:
		return ""
	}
}

// Handles the command to kick a user from the server
func handleBotCommandKick(user *sessions.User, args []string) string {
	if !common.HasPrivilege(user.Info.Privileges, common.PrivilegeKickUsers) {
		return ""
	}

	if len(args) < 2 {
		return "You must specify a user to kick."
	}

	target := getUserFromCommandArgs(args)

	if target == nil {
		return "That user is not online."
	}

	if target == Bot {
		return botMessageStop
	}

	sessions.SendPacketToUser(packets.NewServerNotificationError("You have been kicked from the server."), target)
	utils.CloseConnectionDelayed(target.Conn)
	return fmt.Sprintf("%v has been kicked from the server.", target.Info.Username)
}

// Handles the command to notify all users of something
func handleBotCommandNotifyAll(user *sessions.User, args []string) string {
	if !common.HasPrivilege(user.Info.Privileges, common.PrivilegeNotifyUsers) {
		return ""
	}

	if len(args) < 2 {
		return "You must provide a message to notify everyone with."
	}

	notification := strings.Join(args[1:], " ")

	for _, onlineUser := range sessions.GetOnlineUsers() {
		sessions.SendPacketToUser(packets.NewServerNotificationInfo(notification), onlineUser)
		sendPrivateMessage(Bot, onlineUser, notification)
	}

	return "Your message has been notified to all online users."
}

// getUserFromCommandArgs Returns a target user from command args
func getUserFromCommandArgs(args []string) *sessions.User {
	return sessions.GetUserByUsername(strings.ToLower(strings.ReplaceAll(args[1], "_", " ")))
}
