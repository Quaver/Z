package multiplayer

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
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
	}

	return fmt.Sprintf("You executed the multiplayer command: %v", args)
}

// Handles the command to kick a user
func handleCommandKickPlayer(user *sessions.User, game *Game, args []string) string {
	if !game.isUserHost(user) {
		return ""
	}

	if len(args) < 3 {
		return "You must provide a username to kick."
	}

	target := sessions.GetUserByUsername(strings.ToLower(strings.ReplaceAll(args[2], "_", " ")))

	if target == nil {
		return "That player is not online."
	}

	if target == user {
		return "You cannot kick yourself from the game."
	}

	if !utils.Includes(game.Data.PlayerIds, target.Info.Id) {
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
	return fmt.Sprintf("The multiplayer game has been changed to %v", game.Data.Name)
}
