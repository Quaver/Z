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

	switch strings.ToLower(args[1]) {
	case "kick":
		return handleCommandKickPlayer(user, game, args)
	case "name":
		return handleCommandChangeName(user, game, args)
	case "host":
		return handleCommandChangeHost(user, game, args)
	case "map":
		return handleCommandChangeMap(user, game, args)
	}

	return fmt.Sprintf("You executed the multiplayer command: %v", args)
}

// Handles the command to kick a user
func handleCommandKickPlayer(user *sessions.User, game *Game, args []string) string {
	message := ""

	game.RunLocked(func() {
		if !game.isUserHost(user) {
			message = ""
			return
		}

		if len(args) < 3 {
			message = "You must provide a username to kick."
			return
		}

		target := getUserFromCommandArgs(args)

		if target == nil {
			message = "That player is not online."
			return
		}

		if target == user {
			message = "You cannot kick yourself from the game."
			return
		}

		if !game.isUserInGame(target) {
			message = "That user is not in the game."
			return
		}

		game.KickPlayer(user, target.Info.Id)
		message = fmt.Sprintf("%v has been successfully kicked from the game.", target.Info.Username)
	})

	return message
}

// Handles the command to change the name of the multiplayer game.
func handleCommandChangeName(user *sessions.User, game *Game, args []string) string {
	message := ""

	game.RunLocked(func() {
		if !game.isUserHost(user) {
			message = ""
			return
		}

		if len(args) < 3 {
			message = "You must provide a new name for the multiplayer game."
			return
		}

		game.ChangeName(user, strings.Join(args[2:], " "))
		message = fmt.Sprintf("The multiplayer game name has been changed to: %v.", game.Data.Name)
	})

	return message
}

// Handles the command to change the host of the game
func handleCommandChangeHost(user *sessions.User, game *Game, args []string) string {
	message := ""

	game.RunLocked(func() {
		if !game.isUserHost(user) {
			message = ""
			return
		}

		if len(args) < 3 {
			message = "You must provide the username of the player to give host to."
			return
		}

		target := getUserFromCommandArgs(args)

		if target == nil {
			message = "That user is not online."
			return
		}

		if !game.isUserInGame(target) {
			message = "That user is not in the game."
			return
		}

		game.SetHost(user, target.Info.Id)
		message = fmt.Sprintf("The host has been transferred to: %v.", target.Info.Username)
	})

	return message
}

// Handles the command to change the multiplayer map
// TODO: Needs difficulty calculator
func handleCommandChangeMap(user *sessions.User, game *Game, args []string) string {
	message := ""

	game.RunLocked(func() {
	})

	message = "Command not implemented"
	return message
}

// Returns a target user from command args
func getUserFromCommandArgs(args []string) *sessions.User {
	return sessions.GetUserByUsername(strings.ToLower(strings.ReplaceAll(args[2], "_", " ")))
}
