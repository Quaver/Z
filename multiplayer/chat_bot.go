package multiplayer

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/objects"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
)

func InitializeChatBot() {
	chat.AddPublicMessageHandler(handleMultiplayerCommands)
	chat.AddPublicMessageHandler(handleJoinMultiplayerChatCommand)
}

// Handles the command to join a multiplayer chat channel.
func handleJoinMultiplayerChatCommand(user *sessions.User, channel *chat.Channel, args []string) string {
	if args[0] != "!joinmpchat" {
		return ""
	}

	if !common.HasPrivilege(user.Info.Privileges, common.PrivilegeEnableTournamentMode) {
		return ""
	}

	if len(args) < 2 {
		return "You must provide a game id."
	}

	id, err := strconv.Atoi(args[1])

	if err != nil {
		return "You must provide a valid number id."
	}

	if game, ok := lobby.games[id]; ok {
		game.chatChannel.AddUser(user)
	}

	return "That multiplayer game does not exist."
}

// Handles commands made for multiplayer
func handleMultiplayerCommands(user *sessions.User, channel *chat.Channel, args []string) string {
	if channel.Type != chat.ChannelTypeMultiplayer {
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
		case "allowmode":
			message = handleCommandModeAllowance(user, game, args, true)
		case "disallowmode":
			message = handleCommandModeAllowance(user, game, args, false)
		case "lnmin":
			message = handleCommandLongNote(user, game, args, false)
		case "lnmax":
			message = handleCommandLongNote(user, game, args, true)
		case "freemod":
			message = handleCommandFreeMod(user, game, objects.MultiplayerGameFreeModRegular)
		case "freerate":
			message = handleCommandFreeMod(user, game, objects.MultiplayerGameFreeModRate)
		case "clearwins":
			message = handleCommandClearWins(user, game)
		case "playerwins":
			message = handleCommandPlayerWins(user, game, args)
		case "referee":
			message = handleCommandReferee(user, game, args)
		case "clearreferee":
			message = handleCommandClearReferee(user, game)
		case "tournament":
			message = handleCommandTournamentMode(user, game)
		case "invite":
			message = handleCommandInvite(user, game, args)
		case "roll":
			message = handleCommandRoll(user)
		case "autohost":
			message = handleCommandAutoHost(user, game)
		case "randmap":
			message = handleCommandRandomMap(user, game)
		case "debug":
			message = handleCommandDebug(user, game)
		case "mods":
			message = handleCommandMods(user, game, args)
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
	return ""
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
	return ""
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
	return ""
}

// Handles the command to change the multiplayer map
func handleCommandChangeMap(user *sessions.User, game *Game, args []string) string {
	if !game.isUserHost(user) {
		return ""
	}

	if len(args) < 3 {
		return "You must provide a map id."
	}

	id, err := strconv.Atoi(args[2])

	if err != nil {
		return "You must provide a valid map id."
	}

	song, err := db.GetSongMapById(id)

	if err != nil {
		if err == sql.ErrNoRows {
			return "That map doesn't exist."
		}

		log.Printf("Error getting map %v from the database - %v\n", id, err)
		return "There was an error while retrieving the map."
	}

	game.changeMapFromDbSong(song)
	return ""
}

// Handles the command to enable/disable host rotation
func handleCommandHostRotation(user *sessions.User, game *Game) string {
	if !game.isUserHost(user) {
		return ""
	}

	game.SetHostRotation(user, !game.Data.IsHostRotation)
	return ""
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
	return ""
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
	return ""
}

// Handles the command to end the match
func handleCommandEndMatch(user *sessions.User, game *Game) string {
	if !game.isUserHost(user) {
		return ""
	}

	if !game.Data.InProgress {
		return "The match is not currently in progress."
	}

	game.EndGame(true)
	return ""
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
	return ""
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
	return ""
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

	return ""
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
	return ""
}

// Handles the command to set an allowed game mode for the game
func handleCommandModeAllowance(user *sessions.User, game *Game, args []string, allowing bool) string {
	if !game.isUserHost(user) {
		return ""
	}

	errorStr := "You must provide either `4k` or `7k`"

	if len(args) < 3 {
		return errorStr
	}

	var mode common.Mode

	mode, err := common.GetModeFromShortHand(args[2])

	if err != nil {
		return "Invalid mode provided"
	}

	if allowing && !utils.Includes(game.Data.FilterAllowedGameModes, mode) {
		game.SetAllowedGameModes(user, append(game.Data.FilterAllowedGameModes, mode))
		return fmt.Sprintf("%v is now allowed in this game.", args[2])
	} else if !allowing && len(game.Data.FilterAllowedGameModes) > 1 {
		game.SetAllowedGameModes(user, utils.Filter(game.Data.FilterAllowedGameModes, func(x common.Mode) bool { return x != mode }))
		return fmt.Sprintf("%v is now disallowed in this game.", args[2])
	}

	if allowing {
		return "This mode is allowed already."
	} else {
		return "You must have at least one allowed mode."
	}
}

// Handles the command to change the long note percentage
func handleCommandLongNote(user *sessions.User, game *Game, args []string, isMax bool) string {
	if !game.isUserHost(user) {
		return ""
	}

	if len(args) < 3 {
		return "You must provide a valid number."
	}

	percentage, err := strconv.Atoi(args[2])

	if err != nil {
		return "You must provide a valid number."
	}

	if !isMax && percentage > game.Data.FilterMaxLongNotePercent {
		return "The minimum long note percentage cannot be above the maxim long note percentage."
	} else if isMax && percentage < game.Data.FilterMinLongNotePercent {
		return "The maximum long note percentage cannot be below the minimum long note percentage."
	}

	if isMax {
		game.SetLongNotePercent(user, game.Data.FilterMinLongNotePercent, percentage)
	} else {
		game.SetLongNotePercent(user, percentage, game.Data.FilterMaxLongNotePercent)
	}

	return ""
}

// Handles enabling/disabling free mod / free rate for the game
func handleCommandFreeMod(user *sessions.User, game *Game, freeModType objects.MultiplayerGameFreeMod) string {
	if !game.isUserHost(user) {
		return ""
	}

	if game.Data.FreeModType&freeModType != 0 {
		game.SetFreeMod(user, game.Data.FreeModType-freeModType)
	} else {
		game.SetFreeMod(user, game.Data.FreeModType|freeModType)
	}

	return ""
}

// Handles the command to clear all players' win counts
func handleCommandClearWins(user *sessions.User, game *Game) string {
	if !game.isUserHost(user) {
		return ""
	}

	for _, playerId := range game.Data.PlayerIds {
		game.SetPlayerWinCount(playerId, 0)
	}

	return "All player win counts have been reset back to zero."
}

// Handles the command to set a specific player's win count
func handleCommandPlayerWins(user *sessions.User, game *Game, args []string) string {
	if !game.isUserHost(user) {
		return ""
	}

	if len(args) < 4 {
		return "Invalid command usage. Try this: `!mp playerwins user_name number`."
	}

	wins, err := strconv.Atoi(args[3])

	if err != nil {
		return "You must supply a valid win count."
	}

	wins = utils.Clamp(wins, 0, 100)

	target := getUserFromCommandArgs(args)

	if target == nil {
		return "That user is not online."
	}

	if !game.isUserInGame(target) {
		return "That user is not in the game."
	}

	game.SetPlayerWinCount(target.Info.Id, wins)
	return fmt.Sprintf("%v's win count has been set to: %v.", target.Info.Username, wins)
}

// Handles the command to appoint a user as referee.
func handleCommandReferee(user *sessions.User, game *Game, args []string) string {
	if !game.isUserHost(user) {
		return ""
	}

	if len(args) < 3 {
		return "You must provide a user to give host to."
	}

	target := getUserFromCommandArgs(args)

	if target == nil {
		return "That user is not online."
	}

	if !game.isUserInGame(target) {
		return "That user is not in the game."
	}

	game.SetReferee(user, target.Info.Id)
	return fmt.Sprintf("%v is now the referee of the game.", target.Info.Username)
}

// Handles the command to clear the referee of the game.
func handleCommandClearReferee(user *sessions.User, game *Game) string {
	if !game.isUserHost(user) {
		return ""
	}

	game.SetReferee(user, -1)
	return "The referee of the game has been cleared."
}

func handleCommandTournamentMode(user *sessions.User, game *Game) string {
	if !game.isUserHost(user) {
		return ""
	}

	if !common.HasPrivilege(user.Info.Privileges, common.PrivilegeEnableTournamentMode) {
		return "You don't have permission to turn on tournament mode."
	}

	game.SetTournamentMode(user, !game.Data.IsTournamentMode)
	return ""
}

// Handles the command to invite a user to the game
func handleCommandInvite(user *sessions.User, game *Game, args []string) string {
	if len(args) < 3 {
		return "You must provide a user to invite."
	}

	target := getUserFromCommandArgs(args)

	if target == nil {
		return "That user is not online."
	}

	if game.isUserInGame(target) {
		return "That user is already in the game."
	}

	game.SendInvite(user, target)
	return ""
}

// Handles the command to roll a random number between 0 and 100
func handleCommandRoll(user *sessions.User) string {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(101)

	return fmt.Sprintf("%v has rolled a: %v.", user.Info.Username, randomNumber)
}

// Enables/Disables AutoHost for the game.
func handleCommandAutoHost(user *sessions.User, game *Game) string {
	if !game.isUserHost(user) {
		return ""
	}

	game.SetAutoHost(user, !game.Data.IsAutoHost)
	return ""
}

// Selects a random map in multiplayer
func handleCommandRandomMap(user *sessions.User, game *Game) string {
	if !game.isUserHost(user) {
		return ""
	}

	game.selectAutohostMap()
	return ""
}

func handleCommandDebug(user *sessions.User, game *Game) string {
	if !common.HasUserGroup(user.Info.UserGroups, common.UserGroupSwan) {
		return ""
	}

	str := fmt.Sprintf("In Progress: %v\n", game.Data.InProgress)
	str += fmt.Sprintf("Disbanded: %v\n", game.isDisbanded)
	str += fmt.Sprintf("Player Count: %v\n", len(game.Data.PlayerIds))
	str += fmt.Sprintf("Map rating: %v\n", game.Data.MapDifficultyRating)
	str += fmt.Sprintf("Players:\n")

	for index, id := range game.Data.PlayerIds {
		user := sessions.GetUserById(id)
		userFound := user != nil

		str += fmt.Sprintf("%v. Id: %v | User Online: %v", index, id, userFound)

		if userFound {
			str += fmt.Sprintf(" | Username: %v\n", user.Info.Username)
		} else {
			str += "\n"
		}
	}

	return str
}

// Handles current map modifiers in multiplayer
func handleCommandMods(user *sessions.User, game *Game, args []string) string {
	if !game.isUserHost(user) {
		return ""
	}

	if !game.Data.IsTournamentMode {
		return ""
	}

	if len(args) < 3 {
		return "Incorrect arguments, usage: !mp mods <mod1,mod2,mod3>"
	}

	modMap := common.GetModStrings()
	var mods common.Mods
	var validatedMods []string

	for _, arg := range strings.Split(args[2], ",") {
		mod, exists := modMap[arg]
		if !exists {
			continue
		}

		if mods&mod != 0 {
			continue
		}

		validatedMods = append(validatedMods, arg)

		mods |= mod
	}

	difficulty := game.findMapDifficultyRatingFromMods(mods)

	game.SetGlobalModifiers(user, mods, difficulty)

	if len(validatedMods) == 0 {
		return "All player modifiers have been reset."
	}

	return fmt.Sprintf("Game modifiers changed to: %v.", strings.Join(validatedMods, ","))
}

// getUserFromCommandArgs Returns a target user from command args
func getUserFromCommandArgs(args []string) *sessions.User {
	return sessions.GetUserByUsername(strings.ToLower(strings.ReplaceAll(args[2], "_", " ")))
}
