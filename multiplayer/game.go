package multiplayer

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/objects"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/scoring"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
	"log"
	"math"
	"time"
)

type Game struct {
	mutex               *utils.Mutex                    // Locks down the game to prevent race conditions
	Data                *objects.MultiplayerGame        // Data about the multiplayer game that is sent in a packet
	Password            string                          // The password for the game. This is different from Data.CreationPassword, as it is hidden from users.
	CreatorId           int                             // The id of the user who created the game
	countdownTimer      *time.Timer                     // Counts down before starting the game
	playersInvited      []int                           // A list of users who have been invited to the game
	playersInMatch      []int                           // A list of users who are currently playing the current match
	playersScreenLoaded []int                           // A list of users whose screens have loaded in-game. The match doesn't start until all players are loaded.
	playersFinished     []int                           // A list of users who have finished playing the map
	playersSkipped      []int                           // A list of players who have skipped the map in multiplayer
	playerScores        map[int]*scoring.ScoreProcessor // Score processors for players in the game
	chatChannel         *chat.Channel                   // The multiplayer chat
}

const (
	maxPlayerCount int = 16 // The maximum amount of players allowed in a game
)

// NewGame Creates a new multiplayer game from a game
func NewGame(gameData *objects.MultiplayerGame, creatorId int) (*Game, error) {
	game := Game{
		mutex:               utils.NewMutex(),
		Data:                gameData,
		CreatorId:           creatorId,
		Password:            gameData.CreationPassword,
		playersInvited:      []int{},
		playersInMatch:      []int{},
		playersScreenLoaded: []int{},
		playersFinished:     []int{},
		playersSkipped:      []int{},
		playerScores:        map[int]*scoring.ScoreProcessor{},
	}

	game.Data.GameId = utils.GenerateRandomString(32)
	game.Data.CreationPassword = ""
	game.Data.SetDefaults()
	game.validateSettings()

	var err error
	game.Data.Id, err = db.InsertMultiplayerGame(game.Data.Name, game.Data.GameId)

	if err != nil {
		return nil, err
	}

	game.chatChannel = chat.AddMultiplayerChannel(game.Data.GameId)
	return &game, nil
}

// RunLocked Runs a function in a locked environment
func (game *Game) RunLocked(f func()) {
	game.mutex.RunLocked(f)
}

// AddPlayer Adds a user to the multiplayer game
func (game *Game) AddPlayer(userId int, password string) {
	user := sessions.GetUserById(userId)

	if user == nil {
		return
	}

	currentGame := GetGameById(user.GetMultiplayerGameId())

	if currentGame != nil {
		currentGame.RemovePlayer(user.Info.Id)
	}

	if len(game.Data.PlayerIds) >= maxPlayerCount {
		sessions.SendPacketToUser(packets.NewServerJoinGameFailed(packets.JoinGameErrorFull), user)
		return
	}

	// Check password in the event that the user wasn't invited or has a swan-bypass.
	if (game.Data.HasPassword && game.Password != password) && !utils.Includes(game.playersInvited, userId) && !common.IsSwan(user.Info.UserGroups) {
		sessions.SendPacketToUser(packets.NewServerJoinGameFailed(packets.JoinGameErrorPassword), user)
		return
	}

	game.Data.PlayerIds = append(game.Data.PlayerIds, user.Info.Id)
	game.Data.PlayerModifiers = append(game.Data.PlayerModifiers, &objects.MultiplayerGamePlayerMods{Id: user.Info.Id})
	game.Data.PlayerWins = append(game.Data.PlayerWins, &objects.MultiplayerGamePlayerWins{Id: user.Info.Id})

	user.SetMultiplayerGameId(game.Data.Id)

	game.sendPacketToPlayers(packets.NewServerUserJoinedGame(user.Info.Id))
	sessions.SendPacketToUser(packets.NewServerJoinGame(game.Data.GameId), user)

	if len(game.Data.PlayerIds) == 1 {
		game.SetHost(nil, user.Info.Id)
	}

	game.chatChannel.AddUser(user)

	RemoveUserFromLobby(user)
	sendLobbyUsersGameInfoPacket(game, true)
}

// RemovePlayer Removes a player from the multiplayer game and disbands the game if necessary
func (game *Game) RemovePlayer(userId int) {
	user := sessions.GetUserById(userId)

	if user != nil {
		user.SetMultiplayerGameId(0)
		game.chatChannel.RemoveUser(user)
	}

	game.Data.PlayerIds = utils.Filter(game.Data.PlayerIds, func(x int) bool { return x != userId })
	game.Data.PlayerModifiers = utils.Filter(game.Data.PlayerModifiers, func(x *objects.MultiplayerGamePlayerMods) bool { return x.Id != userId })
	game.Data.PlayerWins = utils.Filter(game.Data.PlayerWins, func(x *objects.MultiplayerGamePlayerWins) bool { return x.Id != userId })
	game.playersInMatch = utils.Filter(game.playersInMatch, func(x int) bool { return x != userId })
	game.playersScreenLoaded = utils.Filter(game.playersScreenLoaded, func(x int) bool { return x != userId })
	game.playersFinished = utils.Filter(game.playersFinished, func(x int) bool { return x != userId })
	game.playersSkipped = utils.Filter(game.playersSkipped, func(x int) bool { return x != userId })
	delete(game.playerScores, userId)

	// Disband game since there are no more players left
	if len(game.Data.PlayerIds) == 0 {
		game.EndGame()
		RemoveGameFromLobby(game)
		chat.RemoveMultiplayerChannel(game.Data.GameId)
		return
	}

	game.SetHost(nil, game.Data.PlayerIds[0])
	game.sendPacketToPlayers(packets.NewServerUserLeftGame(userId))
	game.checkScreenLoadedPlayers()
	game.checkAllPlayersSkipped()

	if game.isAllPlayersFinished() {
		game.EndGame()
	}

	sendLobbyUsersGameInfoPacket(game, true)
}

// KickPlayer Kicks a player from the multiplayer game
func (game *Game) KickPlayer(requester *sessions.User, userId int) {
	if !game.isUserHost(requester) || !utils.Includes(game.Data.PlayerIds, userId) {
		return
	}

	game.RemovePlayer(userId)

	user := sessions.GetUserById(userId)

	if user == nil {
		return
	}

	sessions.SendPacketToUser(packets.NewServerGameKicked(), user)
}

// SetHost Sets the host of the game. Set requester to nil if this is meant to be a forced action.
// Otherwise, set the requester if a user is transferring host to another.
func (game *Game) SetHost(requester *sessions.User, userId int) {
	if !game.isUserHost(requester) {
		return
	}

	if !utils.Includes(game.Data.PlayerIds, userId) {
		return
	}

	game.Data.HostId = userId
	game.SetHostSelectingMap(nil, false, false)

	game.sendPacketToPlayers(packets.NewServerGameChangeHost(game.Data.HostId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// ChangeMap Changes the multiplayer map. Non-nil requester checks if they are the host
func (game *Game) ChangeMap(requester *sessions.User, packet *packets.ClientChangeGameMap) {
	if game.Data.InProgress {
		return
	}

	if !game.isUserHost(requester) {
		return
	}

	game.Data.MapMD5 = packet.MD5
	game.Data.MapMD5Alternative = packet.AlternativeMD5
	game.Data.MapId = packet.MapId
	game.Data.MapsetId = packet.MapsetId
	game.Data.MapName = packet.Name
	game.Data.MapGameMode = packet.Mode
	game.Data.MapDifficultyRating = packet.DifficultyRating
	game.Data.MapDifficultyRatingAll = packet.DifficultyRatingAll
	game.Data.MapJudgementCount = packet.JudgementCount
	game.Data.PlayersWithoutMap = []int{}
	game.Data.PlayersReady = []int{}
	game.clearReadyPlayers(false)
	game.clearCountdown()
	game.validateSettings()

	game.sendPacketToPlayers(packets.NewServerGameMapChanged(packet))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerDoesntHaveMap Sets that a player does not have the map downloaded
func (game *Game) SetPlayerDoesntHaveMap(userId int) {
	game.Data.PlayersWithoutMap = append(game.Data.PlayersWithoutMap, userId)

	game.sendPacketToPlayers(packets.NewServerGamePlayerNoMap(userId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerHasMap Sets that a player now has the currently played map
func (game *Game) SetPlayerHasMap(userId int) {
	game.Data.PlayersWithoutMap = utils.Filter(game.Data.PlayersWithoutMap, func(x int) bool {
		return x != userId
	})

	game.sendPacketToPlayers(packets.NewServerGamePlayerHasMap(userId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerReady Sets that a player is currently ready to play
func (game *Game) SetPlayerReady(userId int) {
	if game.Data.InProgress {
		return
	}

	if !utils.Includes(game.Data.PlayersReady, userId) {
		game.Data.PlayersReady = append(game.Data.PlayersReady, userId)
	}

	game.sendPacketToPlayers(packets.NewServerGamePlayerReady(userId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerNotReady Sets that a player is not ready to play
func (game *Game) SetPlayerNotReady(userId int) {
	game.Data.PlayersReady = utils.Filter(game.Data.PlayersReady, func(i int) bool {
		return i != userId
	})

	game.sendPacketToPlayers(packets.NewServerGamePlayerNotReady(userId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// StartCountdown Starts the 5-second multiplayer countdown
func (game *Game) StartCountdown(requester *sessions.User) {
	if game.Data.InProgress {
		return
	}

	if !game.isUserHost(requester) {
		return
	}

	game.countdownTimer = time.AfterFunc(5*time.Second, func() {
		game.RunLocked(func() {
			game.StartGame()
		})
	})

	game.sendPacketToPlayers(packets.NewServerGameStartCountdown())
	sendLobbyUsersGameInfoPacket(game, true)
}

// StopCountdown Stops the multiplayer countdown if one is live
func (game *Game) StopCountdown(requester *sessions.User) {
	if game.Data.InProgress {
		return
	}

	if !game.isUserHost(requester) {
		return
	}

	game.clearCountdown()
	sendLobbyUsersGameInfoPacket(game, true)
}

// StartGame Starts the multiplayer game
func (game *Game) StartGame() {
	if game.Data.InProgress {
		return
	}

	game.Data.InProgress = true

	game.playersInMatch = utils.Filter(game.Data.PlayerIds, func(x int) bool {
		return x != game.Data.RefereeId && !utils.Includes(game.Data.PlayersWithoutMap, x)
	})

	game.createScoreProcessors()
	game.clearCountdown()
	game.clearReadyPlayers(false)
	game.SetHostSelectingMap(nil, false, false)

	game.sendPacketToPlayers(packets.NewServerGameStart())
	sendLobbyUsersGameInfoPacket(game, true)
}

// EndGame Ends the multiplayer game
func (game *Game) EndGame() {
	if !game.Data.InProgress {
		return
	}

	game.Data.InProgress = false
	game.playersInMatch = []int{}
	game.playersScreenLoaded = []int{}
	game.playersFinished = []int{}
	game.playersSkipped = []int{}
	game.clearCountdown()
	game.clearReadyPlayers(false)
	game.rotateHost()

	game.sendPacketToPlayers(packets.NewServerGameEnded())
	sendLobbyUsersGameInfoPacket(game, true)
}

// ChangeName Changes the name of the multiplayer game
func (game *Game) ChangeName(requester *sessions.User, name string) {
	if game.Data.InProgress {
		return
	}

	if !game.isUserHost(requester) {
		return
	}

	if game.Data.Name == "" {
		return
	}

	game.Data.Name = name
	game.validateSettings()

	game.sendPacketToPlayers(packets.NewServerGameNameChanged(game.Data.Name))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetHostSelectingMap Sets whether the host is selecting a map or not
func (game *Game) SetHostSelectingMap(requester *sessions.User, isSelecting bool, sendToLobby bool) {
	if game.Data.InProgress {
		return
	}

	if !game.isUserHost(requester) {
		return
	}

	game.Data.IsHostSelectingMap = isSelecting
	game.sendPacketToPlayers(packets.NewServerGameHostSelectingMap(isSelecting))

	if sendToLobby {
		sendLobbyUsersGameInfoPacket(game, true)
	}
}

// SetPassword Sets the password for the game
func (game *Game) SetPassword(requester *sessions.User, password string) {
	if !game.isUserHost(requester) {
		return
	}

	game.Password = password
	game.validateSettings()

	sendLobbyUsersGameInfoPacket(game, true)
}

// SetDifficultyRange Sets the difficulty range filter for the game
func (game *Game) SetDifficultyRange(requester *sessions.User, min float32, max float32) {
	if !game.isUserHost(requester) {
		return
	}

	game.Data.FilterMinDifficultyRating = min
	game.Data.FilterMaxDifficultyRating = max
	game.validateSettings()

	packet := packets.NewServerGameDifficultyRangeChanged(game.Data.FilterMinDifficultyRating, game.Data.FilterMaxDifficultyRating)
	game.sendPacketToPlayers(packet)

	sendLobbyUsersGameInfoPacket(game, true)
}

// SetMaxSongLength Sets the maximum song length filter for the map
func (game *Game) SetMaxSongLength(requester *sessions.User, lengthSeconds int) {
	if !game.isUserHost(requester) {
		return
	}

	game.Data.FilterMaxSongLength = lengthSeconds
	game.validateSettings()

	game.sendPacketToPlayers(packets.NewServerGameMaxSongLengthChanged(game.Data.FilterMaxSongLength))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetAllowedGameModes Sets the game modes that are allowed to be played in the game
func (game *Game) SetAllowedGameModes(requester *sessions.User, gameModes []common.Mode) {
	if !game.isUserHost(requester) {
		return
	}

	game.Data.FilterAllowedGameModes = gameModes
	game.validateSettings()

	game.sendPacketToPlayers(packets.NewServerGameAllowedModesChanged(game.Data.FilterAllowedGameModes))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetGlobalModifiers Sets the modifiers that all players must use in the game
func (game *Game) SetGlobalModifiers(requester *sessions.User, mods common.Mods, difficultyRating float64) {
	if game.Data.InProgress {
		return
	}

	if !game.isUserHost(requester) {
		return
	}

	game.Data.GlobalModifiers = mods
	game.Data.MapDifficultyRating = difficultyRating
	game.validateSettings()

	game.sendPacketToPlayers(packets.NewServerGameChangeModifiers(game.Data.GlobalModifiers, game.Data.MapDifficultyRating))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetFreeMod Sets the free mod type for the match (free mod / free rate)
func (game *Game) SetFreeMod(requester *sessions.User, freeMod objects.MultiplayerGameFreeMod) {
	if game.Data.InProgress {
		return
	}

	if !game.isUserHost(requester) {
		return
	}

	game.Data.FreeModType = freeMod
	game.resetAllModifiers()
	game.validateSettings()

	game.sendPacketToPlayers(packets.NewServerGameChangeFreeMod(game.Data.FreeModType))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerModifiers Sets the player modifiers for an individual user
func (game *Game) SetPlayerModifiers(userId int, mods common.Mods) {
	if game.Data.InProgress {
		return
	}

	playerMods, err := utils.Find(game.Data.PlayerModifiers, func(x *objects.MultiplayerGamePlayerMods) bool {
		return x.Id == userId
	})

	if err != nil {
		log.Printf("[MP #v] Error getting playermods for user: #%v - %v\n", userId, err)
		return
	}

	playerMods.Modifiers = mods
	game.sendPacketToPlayers(packets.NewServerGameChangePlayerModifiers(userId, mods))

	sendLobbyUsersGameInfoPacket(game, true)
}

// SetHostRotation Sets whether host rotation will be enabled for the game
func (game *Game) SetHostRotation(requester *sessions.User, enabled bool) {
	if !game.isUserHost(requester) {
		return
	}

	game.Data.IsHostRotation = enabled
	game.sendPacketToPlayers(packets.NewServerGameHostRotation(game.Data.IsHostRotation))

	sendLobbyUsersGameInfoPacket(game, true)
}

// SetLongNotePercent Sets the minimum and maximum long note percentage filters for the game
func (game *Game) SetLongNotePercent(requester *sessions.User, min int, max int) {
	if !game.isUserHost(requester) {
		return
	}

	game.Data.FilterMinLongNotePercent = min
	game.Data.FilterMaxLongNotePercent = max
	game.validateSettings()

	game.sendPacketToPlayers(packets.NewServerGameLongNotePercent(game.Data.FilterMinLongNotePercent, game.Data.FilterMaxLongNotePercent))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetMaxPlayerCount Sets the amount of max players allowed in the game
func (game *Game) SetMaxPlayerCount(requester *sessions.User, count int) {
	if !game.isUserHost(requester) {
		return
	}

	// Can't change max players if there are more players in the game than the requested count
	if len(game.Data.PlayerIds) > count {
		return
	}

	game.Data.MaxPlayers = count
	game.validateSettings()

	game.sendPacketToPlayers(packets.NewServerGameChangeMaxPlayers(game.Data.MaxPlayers))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SendInvite Sends an invitation to a user in the multiplayer game
func (game *Game) SendInvite(sender *sessions.User, user *sessions.User) {
	if user == nil {
		return
	}

	if !utils.Includes(game.playersInvited, user.Info.Id) {
		game.playersInvited = append(game.playersInvited, user.Info.Id)
	}

	sessions.SendPacketToUser(packets.NewServerGameInvite(game.Data.GameId, sender.Info.Username), user)
}

// SetPlayerWinCount Sets the win count for a given player
func (game *Game) SetPlayerWinCount(userId int, wins int) {
	playerWins, err := utils.Find(game.Data.PlayerWins, func(x *objects.MultiplayerGamePlayerWins) bool {
		return x.Id == userId
	})

	if err != nil {
		return
	}

	playerWins.Wins = wins

	game.sendPacketToPlayers(packets.NewServerGamePlayerWinCount(userId, playerWins.Wins))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetReferee Sets the referee for the game
func (game *Game) SetReferee(requester *sessions.User, userId int) {
	if !game.isUserHost(requester) {
		return
	}

	if userId != -1 && !utils.Includes(game.Data.PlayerIds, userId) {
		return
	}

	game.Data.RefereeId = userId

	game.sendPacketToPlayers(packets.NewServerGameSetReferee(game.Data.RefereeId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerScreenLoaded Handles when a client states that their gameplay screen has loaded at the start of a match
func (game *Game) SetPlayerScreenLoaded(userId int) {
	if !game.Data.InProgress || !utils.Includes(game.playersInMatch, userId) {
		return
	}

	if !utils.Includes(game.playersScreenLoaded, userId) {
		game.playersScreenLoaded = append(game.playersScreenLoaded, userId)
	}

	game.checkScreenLoadedPlayers()
}

// SetPlayerFinished Handles when a client states that they have finished playing the current match
func (game *Game) SetPlayerFinished(userId int) {
	if !game.Data.InProgress || !utils.Includes(game.playersInMatch, userId) {
		return
	}

	if !utils.Includes(game.playersFinished, userId) {
		game.playersFinished = append(game.playersFinished, userId)
	}

	if game.isAllPlayersFinished() {
		game.EndGame()
	}
}

// SetPlayerSkippedSong Handles when the client requests to skip the song
func (game *Game) SetPlayerSkippedSong(userId int) {
	if !game.Data.InProgress || !utils.Includes(game.playersInMatch, userId) {
		return
	}

	if !utils.Includes(game.playersSkipped, userId) {
		game.playersSkipped = append(game.playersSkipped, userId)
	}

	game.checkAllPlayersSkipped()
}

// HandlePlayerJudgements Handles when a player sends judgement data during a multiplayer match
func (game *Game) HandlePlayerJudgements(userId int, judgements []common.Judgements) {
	if !game.Data.InProgress || !utils.Includes(game.playersInMatch, userId) {
		return
	}

	if score, ok := game.playerScores[userId]; ok {
		score.AddJudgements(judgements)
	}

	packet := packets.NewServerGameJudgements(userId, judgements)

	for _, playerId := range game.playersInMatch {
		if playerId == userId {
			continue
		}

		player := sessions.GetUserById(playerId)

		if player == nil {
			continue
		}

		sessions.SendPacketToUser(packet, player)
	}
}

// SetTournamentMode Enables/disables tournament mode for the match
func (game *Game) SetTournamentMode(requester *sessions.User, enabled bool) {
	if !game.isUserHost(requester) {
		return
	}

	if requester != nil && !common.HasPrivilege(requester.Info.Privileges, common.PrivilegeEnableTournamentMode) {
		return
	}

	game.Data.IsTournamentMode = enabled

	game.sendPacketToPlayers(packets.NewServerGameTournamentMode(game.Data.IsTournamentMode))
	sendLobbyUsersGameInfoPacket(game, true)
}

// rotateHost Rotates the host to the next person in line.
func (game *Game) rotateHost() {
	if !game.Data.IsHostRotation {
		return
	}

	if len(game.Data.PlayerIds) == 1 {
		return
	}

	index := utils.FindIndex(game.Data.PlayerIds, game.Data.HostId)

	if index == -1 {
		return
	}

	// Cyclically rotates the host
	if index+1 < len(game.Data.PlayerIds) {
		game.SetHost(nil, game.Data.PlayerIds[index+1])
	} else {
		game.SetHost(nil, game.Data.PlayerIds[0])
	}
}

// Returns if the user is host of the game or has permission.
func (game *Game) isUserHost(user *sessions.User) bool {
	if user == nil {
		return true
	}

	if user.Info.Id != game.Data.HostId {
		return false
	}

	return true
}

// Clears all players that are ready.
func (game *Game) clearReadyPlayers(sendToLobby bool) {
	for _, id := range game.Data.PlayersReady {
		game.sendPacketToPlayers(packets.NewServerGamePlayerNotReady(id))
	}

	game.Data.PlayersReady = []int{}

	if sendToLobby {
		sendLobbyUsersGameInfoPacket(game, true)
	}
}

// Creates score processors for all the users that are playing in the match
func (game *Game) createScoreProcessors() {
	for _, player := range game.playersInMatch {
		playerMods, err := utils.Find(game.Data.PlayerModifiers, func(x *objects.MultiplayerGamePlayerMods) bool {
			return x.Id == player
		})

		if err != nil {
			playerMods = &objects.MultiplayerGamePlayerMods{
				Id:        player,
				Modifiers: 0,
			}
		}

		mods := game.Data.GlobalModifiers | playerMods.Modifiers
		difficulty := game.Data.MapDifficultyRating
		idx := utils.FindIndex(common.SpeedMods, common.GetSpeedModFromMods(mods))

		if idx != -1 && len(game.Data.MapDifficultyRatingAll) > 0 {
			difficulty = game.Data.MapDifficultyRatingAll[idx]
		}

		game.playerScores[player] = scoring.NewScoreProcessor(difficulty, mods)
	}
}

// Clears and stops the countdown timer.
func (game *Game) clearCountdown() {
	if game.countdownTimer != nil {
		game.countdownTimer.Stop()
		game.countdownTimer = nil
	}

	game.sendPacketToPlayers(packets.NewServerGameStopCountdown())
}

// Resets the modifiers for every player
func (game *Game) resetAllModifiers() {
	game.Data.GlobalModifiers = 0
	game.sendPacketToPlayers(packets.NewServerGameChangeModifiers(0, game.Data.MapDifficultyRating))

	for _, pm := range game.Data.PlayerModifiers {
		pm.Modifiers = 0
		game.sendPacketToPlayers(packets.NewServerGameChangePlayerModifiers(pm.Id, pm.Modifiers))
	}
}

// Performs a check if all the players that are playing have their screens loaded, then sends a packet.
// Splitting this out into its own function because we need to check this in multiple places -
// such as when a player leaves a match before their screen loads. It'll prevent it from getting stuck
func (game *Game) checkScreenLoadedPlayers() {
	if !game.Data.InProgress {
		return
	}

	for _, player := range game.playersInMatch {
		if !utils.Includes(game.playersScreenLoaded, player) {
			return
		}
	}

	game.sendPacketToPlayers(packets.NewServerGameAllPlayersLoaded())
}

// Returns if all users that are playing have finished playing the map
func (game *Game) isAllPlayersFinished() bool {
	if !game.Data.InProgress {
		return false
	}

	for _, player := range game.playersInMatch {
		if !utils.Includes(game.playersFinished, player) {
			return false
		}
	}

	return true
}

// Checks if all the players in the game have skipped the map and sends a packet letting them know.
func (game *Game) checkAllPlayersSkipped() {
	if !game.Data.InProgress {
		return
	}

	for _, player := range game.playersInMatch {
		if !utils.Includes(game.playersSkipped, player) {
			return
		}
	}

	game.sendPacketToPlayers(packets.NewServerGameAllPlayersSkipped())
}

// Sends a packet to all players in the game.
func (game *Game) sendPacketToPlayers(packet interface{}) {
	for _, id := range game.Data.PlayerIds {
		user := sessions.GetUserById(id)

		if user == nil {
			continue
		}

		sessions.SendPacketToUser(packet, user)
	}
}

// validateSettings Checks the multiplayer settings to see if they are in an acceptable range
func (game *Game) validateSettings() {
	data := game.Data

	data.Name = utils.TruncateString(data.Name, 50)
	data.HasPassword = game.Password != ""
	data.MaxPlayers = utils.Clamp(data.MaxPlayers, 2, 16)
	data.Ruleset = objects.MultiplayerGameRulesetFreeForAll
	data.FreeModType = utils.Clamp(data.FreeModType, objects.MultiplayerGameFreeModNone, objects.MultiplayerGameFreeModRegular|objects.MultiplayerGameFreeModRate)

	data.MapMD5 = utils.TruncateString(data.MapMD5, 64)
	data.MapMD5Alternative = utils.TruncateString(data.MapMD5Alternative, 64)
	data.MapName = utils.TruncateString(data.MapName, 250)
	data.MapGameMode = utils.Clamp(data.MapGameMode, common.ModeKeys4, common.ModeKeys7)

	data.FilterMinDifficultyRating = utils.Clamp(data.FilterMinDifficultyRating, 0, math.MaxInt32)
	data.FilterMaxDifficultyRating = utils.Clamp(data.FilterMaxDifficultyRating, 0, math.MaxInt32)
	data.FilterMaxSongLength = utils.Clamp(data.FilterMaxSongLength, 0, math.MaxInt32)
	data.FilterMinLongNotePercent = utils.Clamp(data.FilterMinLongNotePercent, 0, 100)
	data.FilterMaxLongNotePercent = utils.Clamp(data.FilterMaxLongNotePercent, 0, 100)
	data.FilterMinAudioRate = utils.Clamp(data.FilterMinAudioRate, 0.5, 2.0)

	// There is a maximum of 31 rates allowed in the game. So if we don't have all of them, then just clear it.
	if len(data.MapDifficultyRatingAll) < 31 {
		data.MapDifficultyRatingAll = []float64{}
	}

	if len(data.FilterAllowedGameModes) == 0 {
		data.FilterAllowedGameModes = []common.Mode{common.ModeKeys4, common.ModeKeys7}
	}
}
