package multiplayer

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/sessions"
	"fmt"
	"strings"
)

func InitializeChatBot() {
	chat.AddPublicMessageHandler(handleMultiplayerCommands)
}

// Handles commands made for multiplayer
func handleMultiplayerCommands(user *sessions.User, channel *chat.Channel, args []string) string {
	if !channel.IsMultiplayer {
		return ""
	}

	if len(args) < 2 || args[0] != "!mp" {
		return ""
	}

	game := GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return ""
	}

	message := ""

	game.RunLocked(func() {
		switch strings.ToLower(args[1]) {
		case "kick":
			message = handleCommandKickPlayer(user, game, args)
		case "name":
			message = handleCommandChangeName(user, game, args)
		case "host":
			message = handleCommandChangeHost(user, game, args)
		case "map":
			message = handleCommandChangeMap(user, game, args)
		case "hostrotation":
			message = handleCommandHostRotation(user, game)
		}
	})

	return message
}

// Handles the command to kick a user
func handleCommandKickPlayer(user *sessions.User, game *Game, args []string) string {
	if !game.isUserHost(user) {
		return ""
	}

	if len(args) < 3 {
		return "You must provide a username to kick."
	}

	target := getUserFromCommandArgs(args)

	if target == nil {
		return "That player is not online."
	}

	if target == user {
		return "You cannot kick yourself from the game."
	}

	if !game.isUserInGame(target) {
		return "That user is not in the game."
	}

	game.KickPlayer(user, target.Info.Id)
	return fmt.Sprintf("%v has been successfully kicked from the game.", target.Info.Username)
}

// Handles the command to change the name of the multiplayer game.
func handleCommandChangeName(user *sessions.User, game *Game, args []string) string {
	if !game.isUserHost(user) {
		return ""
	}

	if len(args) < 3 {
		return "You must provide a new name for the multiplayer game."
	}

	game.ChangeName(user, strings.Join(args[2:], " "))
	return fmt.Sprintf("The multiplayer game name has been changed to: %v.", game.Data.Name)
}

// Handles the command to change the host of the game
func handleCommandChangeHost(user *sessions.User, game *Game, args []string) string {
	if !game.isUserHost(user) {
		return ""
	}

	if len(args) < 3 {
		return "You must provide the username of the player to give host to."
	}

	target := getUserFromCommandArgs(args)

	if target == nil {
		return "That user is not online."
	}

	if !game.isUserInGame(target) {
		return "That user is not in the game."
	}

	game.SetHost(user, target.Info.Id)
	return fmt.Sprintf("The host has been transferred to: %v.", target.Info.Username)
}

// TODO: Handles the command to change the multiplayer map / Needs difficulty calculator
func handleCommandChangeMap(user *sessions.User, game *Game, args []string) string {
	return "Command not implemented"
}

// Handles the command to enable/disable host rotation
func handleCommandHostRotation(user *sessions.User, game *Game) string {
	if !game.isUserHost(user) {
		return ""
	}

	game.SetHostRotation(user, !game.Data.IsHostRotation)
	return fmt.Sprintf("Host Rotation has been set to: %v.", game.Data.IsHostRotation)
}

// Returns a target user from command args
func getUserFromCommandArgs(args []string) *sessions.User {
	return sessions.GetUserByUsername(strings.ToLower(strings.ReplaceAll(args[2], "_", " ")))
}
