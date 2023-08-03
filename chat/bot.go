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
	"strconv"
	"strings"
	"time"
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
		AvatarUrl:      sql.NullString{String: "1"},
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
	case "alertuser", "notifyuser":
		return handleBotCommandNotifyUser(user, args)
	case "mute":
		return handleBotCommandMuteUser(user, args)
	case "unmute":
		return handleBotCommandUnmuteUser(user, args)
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

// Handles the command to notify a specific user of something.
func handleBotCommandNotifyUser(user *sessions.User, args []string) string {
	if !common.HasPrivilege(user.Info.Privileges, common.PrivilegeNotifyUsers) {
		return ""
	}

	if len(args) < 2 {
		return "You must specify a user to send a notification to."
	}

	target := getUserFromCommandArgs(args)

	if target == nil {
		return "That user is not online."
	}

	if target == Bot {
		return "I'm a bot. I don't need to receive notifications."
	}

	if len(args) < 3 {
		return "You must provide a message to notify this user with."
	}

	notification := strings.Join(args[2:], " ")

	sessions.SendPacketToUser(packets.NewServerNotificationInfo(notification), target)
	sendPrivateMessage(Bot, target, notification)

	return fmt.Sprintf("Your notification has been sent to: %v.", target.Info.Username)
}

// Handles the command to mute a given user
func handleBotCommandMuteUser(user *sessions.User, args []string) string {
	if !common.HasPrivilege(user.Info.Privileges, common.PrivilegeMuteUsers) {
		return ""
	}

	if len(args) < 2 {
		return "You must specify a user to mute."
	}

	target, err := db.GetUserByUsername(strings.ToLower(strings.ReplaceAll(args[1], "_", " ")))

	if err != nil {
		if err == sql.ErrNoRows {
			return "That user does not exist."
		}

		log.Printf("Error retrieving user from the database - %v\n", err)
		return "An error occurred while executing this command."
	}

	if common.HasUserGroup(target.UserGroups, common.UserGroupBot) {
		return "You cannot mute bot users."
	}

	if len(args) < 3 {
		return "You must provide a time value."
	}

	timeVal, err := strconv.Atoi(args[2])

	if err != nil {
		return "You must provide a valid number as a timeVal value."
	}

	if len(args) < 4 {
		return "You must provide a duration value (ex. s/m/h/d)."
	}

	var duration time.Duration

	switch strings.ToLower(args[3]) {
	case "s":
		duration = time.Second * time.Duration(timeVal)
	case "m":
		duration = time.Minute * time.Duration(timeVal)
	case "h":
		duration = time.Hour * time.Duration(timeVal)
	case "d":
		duration = time.Hour * 24 * time.Duration(timeVal)
	default:
		return "You have specified an invalid duration of time."
	}

	// Try to find the user online. If not, create a "temp" user since that has mute capabilities.
	// Doing it this way because we need to update the mute time of the online user right away,
	// but if they're offline, we can just skip to updating it in the DB.
	onlineUser := getUserFromCommandArgs(args)

	if onlineUser != nil {
		err = onlineUser.MuteUser(duration)
	} else {
		err = sessions.NewUser(nil, target).MuteUser(duration)
	}

	if err != nil {
		log.Printf("Error while muting user - %v - %v\n", target.Id, err)
		return "An error occurred while muting this user."
	}

	if duration == time.Duration(0) {
		return fmt.Sprintf("%v has been unmuted.", target.Username)
	}

	return fmt.Sprintf("%v has been muted for %v", target.Username, duration.String())
}

// Handles the command to unmute a user.
func handleBotCommandUnmuteUser(user *sessions.User, args []string) string {
	if len(args) < 2 {
		return "You must specify a user to unmute."
	}

	// Simply mute the user for zero seconds to unmute them.
	return handleBotCommandMuteUser(user, []string{"unmute", args[1], "0", "s"})
}

// getUserFromCommandArgs Returns a target user from command args
func getUserFromCommandArgs(args []string) *sessions.User {
	return sessions.GetUserByUsername(strings.ToLower(strings.ReplaceAll(args[1], "_", " ")))
}
