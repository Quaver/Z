package multiplayer

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
	"fmt"
	"strconv"
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
		case "maxplayers":
			message = handleCommandMaxPlayers(user, game, args)
		case "start":
			message = handleCommandStartMatch(user, game)
		case "end":
			message = handleCommandEndMatch(user, game)
		case "startcountdown":
			message = handleCommandStartCountdown(user, game)
		case "stopcountdown":
			message = handleCommandStopCountdown(user, game)
		case "mindiff":
			message = handleCommandDifficulty(user, game, args, false)
		case "maxdiff":
			message = handleCommandDifficulty(user, game, args, true)
		case "maxlength":
			message = handleCommandMaxLength(user, game, args)
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
	return fmt.Sprintf("Host Rotation has been %v.", utils.BoolToEnabledString(game.Data.IsHostRotation))
}

// Handles the command to set the max player count
func handleCommandMaxPlayers(user *sessions.User, game *Game, args []string) string {
	if !game.isUserHost(user) {
		return ""
	}

	if len(args) < 3 {
		return "You must provide a number between 2 and 16 in order to change the max player count."
	}

	numPlayers, err := strconv.Atoi(args[2])

	if err != nil {
		return "You must provide a valid number."
	}

	game.SetMaxPlayerCount(user, numPlayers)
	return fmt.Sprintf("The max player count has been changed to: %v.", game.Data.MaxPlayers)
}

// Handles the command to start the match
func handleCommandStartMatch(user *sessions.User, game *Game) string {
	if !game.isUserHost(user) {
		return ""
	}

	if game.Data.InProgress {
		return "The match is already in progress."
	}

	game.StartGame()
	return "The match has been started."
}

// Handles the command to end the match
func handleCommandEndMatch(user *sessions.User, game *Game) string {
	if !game.isUserHost(user) {
		return ""
	}

	if !game.Data.InProgress {
		return "The match is not currently in progress."
	}

	game.EndGame()
	return "The match has been ended."
}

// Handles the command to start the match countdown.
func handleCommandStartCountdown(user *sessions.User, game *Game) string {
	if !game.isUserHost(user) {
		return ""
	}

	if game.Data.InProgress {
		return "The match is currently in progress."
	}

	if game.countdownTimer != nil {
		return "The countdown is already active."
	}

	game.StartCountdown(user)
	return "Countdown active. The match will begin in 5 seconds."
}

// Handles the command to stop the match countdown.
func handleCommandStopCountdown(user *sessions.User, game *Game) string {
	if !game.isUserHost(user) {
		return ""
	}

	if game.Data.InProgress {
		return "The match is currently in progress."
	}

	if game.countdownTimer == nil {
		return "The countdown is not currently active."
	}

	game.StopCountdown(user)
	return "The match countdown has been disabled."
}

// Handles the command to set the minimum/maximum difficulty
func handleCommandDifficulty(user *sessions.User, game *Game, args []string, isMax bool) string {
	if !game.isUserHost(user) {
		return ""
	}

	if len(args) < 3 {
		return "You must provide a difficulty number."
	}

	diff, err := strconv.ParseFloat(args[2], 32)

	if err != nil {
		return "You must provide a valid number."
	}

	diffFloat32 := float32(diff)

	if !isMax && diffFloat32 > game.Data.FilterMaxDifficultyRating {
		return "The minimum difficulty rating cannot be above the maximum difficulty rating."
	} else if isMax && diffFloat32 < game.Data.FilterMinDifficultyRating {
		return "The maximum difficulty rating cannot be below the minimum difficulty rating."
	}

	if isMax {
		game.SetDifficultyRange(user, game.Data.FilterMinDifficultyRating, diffFloat32)
	} else {
		game.SetDifficultyRange(user, diffFloat32, game.Data.FilterMaxDifficultyRating)
	}

	return fmt.Sprintf("The difficulty range has been changed to: %v - %v.", game.Data.FilterMinDifficultyRating, game.Data.FilterMaxDifficultyRating)
}

// Handles the command to set the max length in the multiplayer game
func handleCommandMaxLength(user *sessions.User, game *Game, args []string) string {
	if !game.isUserHost(user) {
		return ""
	}

	if len(args) < 3 {
		return "You must provide a value in seconds."
	}

	seconds, err := strconv.Atoi(args[2])

	if err != nil {
		return "You must provide a valid number."
	}

	game.SetMaxSongLength(user, seconds)
	return fmt.Sprintf("The maximum song length has been changed to: %v seconds.", game.Data.FilterMaxSongLength)
}

// Returns a target user from command args
func getUserFromCommandArgs(args []string) *sessions.User {
	return sessions.GetUserByUsername(strings.ToLower(strings.ReplaceAll(args[2], "_", " ")))
}
