package multiplayer

import (
	"errors"
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/objects"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/scoring"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
	"fmt"
	"log"
	"math"
	"slices"
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
	spectators          []int                           // The players who are currently spectating the game
	isDisbanded         bool                            // If the game has been disbanded
}

const (
	countDifficultyRatings int = 31 // The amount of difficulty ratings needed for a map (31 different rates)
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
		spectators:          []int{},
	}

	game.Data.GameId = utils.GenerateRandomString(32)
	game.Data.CreationPassword = ""
	game.Data.SetDefaults()

	var err error
	game.Data.Id, err = db.InsertMultiplayerGame(game.Data.Name, game.Data.GameId)

	if err != nil {
		return nil, err
	}

	game.validateAndCacheSettings()
	game.removeInactivePlayers()

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

	if len(game.Data.PlayerIds) >= game.Data.MaxPlayers {
		sessions.SendPacketToUser(packets.NewServerJoinGameFailed(packets.JoinGameErrorFull), user)
		return
	}

	// Check password in the event that the user wasn't invited or has a swan-bypass.
	if (game.Data.HasPassword && game.Password != password) && !utils.Includes(game.playersInvited, userId) && !common.IsSwan(user.Info.UserGroups) {
		sessions.SendPacketToUser(packets.NewServerJoinGameFailed(packets.JoinGameErrorPassword), user)
		return
	}

	if utils.Includes(game.playersInvited, userId) {
		game.playersInvited = utils.Filter(game.playersInvited, func(x int) bool { return x != userId })
	}

	game.Data.PlayerIds = append(game.Data.PlayerIds, user.Info.Id)
	game.Data.PlayerModifiers = append(game.Data.PlayerModifiers, &objects.MultiplayerGamePlayerMods{Id: user.Info.Id})

	// Player wins persist even if a user leaves and joins the game pack at a later time
	_, err := utils.Find(game.Data.PlayerWins, func(x *objects.MultiplayerGamePlayerWins) bool { return x.Id == user.Info.Id })
	if err != nil {
		game.Data.PlayerWins = append(game.Data.PlayerWins, &objects.MultiplayerGamePlayerWins{Id: user.Info.Id})
	}

	user.SetMultiplayerGameId(game.Data.Id)
	user.StopSpectatingAll()

	game.cachePlayer(user.Info.Id)
	game.chatChannel.AddUser(user)

	if len(game.Data.PlayerIds) == 1 {
		game.SetHost(nil, user.Info.Id)
	}

	if game.Data.RefereeId == user.Info.Id {
		game.SetReferee(nil, user.Info.Id)
	}

	RemoveUserFromLobby(user)

	game.sendBotMessage(fmt.Sprintf("%v has joined the game.", user.Info.Username))

	sessions.SendPacketToUser(packets.NewServerMultiplayerGameInfo(game.Data), user)
	sessions.SendPacketToUser(packets.NewServerJoinGame(game.Data.GameId), user)
	game.sendPacketToPlayers(packets.NewServerUserJoinedGame(user.Info.Id))
	sendLobbyUsersGameInfoPacket(game, true)
}

// RemovePlayer Removes a player from the multiplayer game and disbands the game if necessary
func (game *Game) RemovePlayer(userId int) {
	user := sessions.GetUserById(userId)

	var playerWasInMatch = slices.Contains(game.playersInMatch, userId)

	if user != nil {
		user.SetMultiplayerGameId(0)
		user.StopSpectatingAll()
		game.chatChannel.RemoveUser(user)
		game.sendBotMessage(fmt.Sprintf("%v has left the game.", user.Info.Username))
	}

	game.Data.PlayerIds = utils.Filter(game.Data.PlayerIds, func(x int) bool { return x != userId })
	game.Data.PlayerModifiers = utils.Filter(game.Data.PlayerModifiers, func(x *objects.MultiplayerGamePlayerMods) bool { return x.Id != userId })
	game.playersInMatch = utils.Filter(game.playersInMatch, func(x int) bool { return x != userId })
	game.playersScreenLoaded = utils.Filter(game.playersScreenLoaded, func(x int) bool { return x != userId })
	game.playersFinished = utils.Filter(game.playersFinished, func(x int) bool { return x != userId })
	game.playersSkipped = utils.Filter(game.playersSkipped, func(x int) bool { return x != userId })
	game.spectators = utils.Filter(game.spectators, func(x int) bool { return x != userId })
	game.deleteCachedPlayer(userId)
	delete(game.playerScores, userId)

	// Disband game since there are no more players left
	if len(game.Data.PlayerIds) == 0 {
		game.disband()
		return
	}

	if game.Data.HostId == userId {
		game.SetHost(nil, game.Data.PlayerIds[0])
	}

	game.sendPacketToPlayers(packets.NewServerUserLeftGame(userId))
	game.checkScreenLoadedPlayers()
	game.checkAllPlayersSkipped()

	// The game ends if everyone finishes the gameplay
	// or if we're in a tournament and someone that is neither a referee or a spectator quit
	if game.isAllPlayersFinished() ||
		game.Data.IsTournamentMode && playerWasInMatch {
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

	game.sendBotMessage(fmt.Sprintf("%v has been kicked from the game.", user.Info.Username))
	sessions.SendPacketToUser(packets.NewServerGameKicked(), user)
}

// AddSpectator Adds a spectator to the game.
func (game *Game) AddSpectator(user *sessions.User, password string) {
	// Require the user to have either Donator or EnableTournamentMode in order to spectate
	if !common.HasUserGroup(user.Info.UserGroups, common.UserGroupDonator) && !common.HasPrivilege(user.Info.Privileges, common.PrivilegeEnableTournamentMode) {
		return
	}

	if (game.Data.HasPassword && game.Password != password) && !common.IsSwan(user.Info.UserGroups) {
		sessions.SendPacketToUser(packets.NewServerJoinGameFailed(packets.JoinGameErrorPassword), user)
		return
	}

	if utils.Includes(game.spectators, user.Info.Id) {
		return
	}

	currGame := GetGameById(user.GetMultiplayerGameId())

	if currGame != nil && currGame != game {
		currGame.RemovePlayer(user.Info.Id)
	}

	game.spectators = append(game.spectators, user.Info.Id)
	game.chatChannel.AddUser(user)
	user.SetMultiplayerGameId(game.Data.Id)
	RemoveUserFromLobby(user)

	game.sendBotMessage(fmt.Sprintf("%v has started spectating the game.", user.Info.Username))
	sessions.SendPacketToUser(packets.NewServerSpectateMultiplayerGame(game.Data.GameId), user)
	sendLobbyUsersGameInfoPacket(game, true)

	if game.Data.InProgress {
		game.initializeSpectator(user)
		sessions.SendPacketToUser(packets.NewServerGameStart(), user)
	}
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

	user := sessions.GetUserById(userId)

	if user != nil {
		game.sendBotMessage(fmt.Sprintf("%v is now the host of the game.", user.Info.Username))
	}

	game.Data.HostId = userId
	game.SetHostSelectingMap(nil, false, false)
	game.validateAndCacheSettings()

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
	game.Data.PlayersWithoutMap = []int{}
	game.Data.PlayersReady = []int{}
	game.clearReadyPlayers(false)
	game.clearCountdown()
	game.SetDonatorMapsetShared(false, false)
	game.validateAndCacheSettings()

	game.sendBotMessage(fmt.Sprintf("The map has been changed to: %v.", game.Data.MapName))
	game.sendPacketToPlayers(packets.NewServerGameMapChanged(packet))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerDoesntHaveMap Sets that a player does not have the map downloaded
func (game *Game) SetPlayerDoesntHaveMap(userId int) {
	game.Data.PlayersWithoutMap = append(game.Data.PlayersWithoutMap, userId)
	game.cachePlayer(userId)

	game.sendPacketToPlayers(packets.NewServerGamePlayerNoMap(userId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerHasMap Sets that a player now has the currently played map
func (game *Game) SetPlayerHasMap(userId int) {
	game.Data.PlayersWithoutMap = utils.Filter(game.Data.PlayersWithoutMap, func(x int) bool { return x != userId })
	game.cachePlayer(userId)

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

	game.cachePlayer(userId)

	game.sendPacketToPlayers(packets.NewServerGamePlayerReady(userId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerNotReady Sets that a player is not ready to play
func (game *Game) SetPlayerNotReady(userId int) {
	game.Data.PlayersReady = utils.Filter(game.Data.PlayersReady, func(i int) bool { return i != userId })
	game.cachePlayer(userId)

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

	game.sendBotMessage("The countdown has started. The match will start in 5 seconds.")
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

	game.sendBotMessage("The match countdown has stopped.")
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

	// Force clear replay frames from the server
	for _, playerId := range game.playersInMatch {
		user := sessions.GetUserById(playerId)
		user.ClearReplayFrames()
	}

	game.initializeSpectators()
	game.createScoreProcessors()
	game.clearCountdown()
	game.clearReadyPlayers(false)
	game.SetHostSelectingMap(nil, false, false)
	game.validateAndCacheSettings()

	game.sendBotMessage("The match has been started.")
	game.sendPacketToPlayers(packets.NewServerGameStart())
	sendLobbyUsersGameInfoPacket(game, true)
}

// EndGame Ends the multiplayer game
func (game *Game) EndGame() {
	if !game.Data.InProgress {
		return
	}

	game.clearCountdown()
	game.clearReadyPlayers(false)
	game.updatePlayerWinCount()
	game.insertMatchIntoDatabase()
	game.rotateHost()

	game.Data.InProgress = false
	game.playersInMatch = []int{}
	game.playersScreenLoaded = []int{}
	game.playersFinished = []int{}
	game.playersSkipped = []int{}
	game.playerScores = map[int]*scoring.ScoreProcessor{}

	if game.Data.IsAutoHost {
		game.selectAutohostMap()
	}

	game.validateAndCacheSettings()

	game.sendBotMessage("The match has ended.")
	game.sendPacketToPlayers(packets.NewServerGameEnded())

	for _, spectatorId := range game.spectators {
		var spectator = sessions.GetUserById(spectatorId)
		sessions.SendPacketToUser(packets.NewServerGameEnded(), spectator)
	}

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
	game.validateAndCacheSettings()

	game.sendBotMessage(fmt.Sprintf("The multiplayer game name has been changed to: %v.", game.Data.Name))
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
	game.validateAndCacheSettings()

	sendLobbyUsersGameInfoPacket(game, true)
}

// SetDifficultyRange Sets the difficulty range filter for the game
func (game *Game) SetDifficultyRange(requester *sessions.User, min float32, max float32) {
	if !game.isUserHost(requester) {
		return
	}

	game.Data.FilterMinDifficultyRating = min
	game.Data.FilterMaxDifficultyRating = max
	game.validateAndCacheSettings()

	packet := packets.NewServerGameDifficultyRangeChanged(game.Data.FilterMinDifficultyRating, game.Data.FilterMaxDifficultyRating)
	game.sendPacketToPlayers(packet)

	sendLobbyUsersGameInfoPacket(game, true)
	game.sendBotMessage(fmt.Sprintf("The difficulty range has been changed to: %v - %v.", game.Data.FilterMinDifficultyRating, game.Data.FilterMaxDifficultyRating))
}

// SetMaxSongLength Sets the maximum song length filter for the map
func (game *Game) SetMaxSongLength(requester *sessions.User, lengthSeconds int) {
	if !game.isUserHost(requester) {
		return
	}

	game.Data.FilterMaxSongLength = lengthSeconds
	game.validateAndCacheSettings()

	game.sendBotMessage(fmt.Sprintf("The maximum song length has been changed to: %v seconds.", game.Data.FilterMaxSongLength))
	game.sendPacketToPlayers(packets.NewServerGameMaxSongLengthChanged(game.Data.FilterMaxSongLength))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetAllowedGameModes Sets the game modes that are allowed to be played in the game
func (game *Game) SetAllowedGameModes(requester *sessions.User, gameModes []common.Mode) {
	if !game.isUserHost(requester) {
		return
	}

	game.Data.FilterAllowedGameModes = gameModes
	game.validateAndCacheSettings()

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
	game.validateAndCacheSettings()

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
	game.validateAndCacheSettings()

	game.sendBotMessage("Free Mod type has been changed. All modifiers have been reset.")
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
	game.cachePlayer(userId)

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
	game.validateAndCacheSettings()

	game.sendBotMessage(fmt.Sprintf("Host Rotation has been %v.", utils.BoolToEnabledString(game.Data.IsHostRotation)))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetLongNotePercent Sets the minimum and maximum long note percentage filters for the game
func (game *Game) SetLongNotePercent(requester *sessions.User, min int, max int) {
	if !game.isUserHost(requester) {
		return
	}

	game.Data.FilterMinLongNotePercent = min
	game.Data.FilterMaxLongNotePercent = max
	game.validateAndCacheSettings()

	game.sendBotMessage(fmt.Sprintf("The long note percentage range has been changed to: %v - %v", game.Data.FilterMinLongNotePercent, game.Data.FilterMaxLongNotePercent))
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
	game.validateAndCacheSettings()

	game.sendBotMessage(fmt.Sprintf("The max player count has been changed to: %v.", game.Data.MaxPlayers))
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

	game.sendBotMessage(fmt.Sprintf("%v has invited %v to the game.", sender.Info.Username, user.Info.Username))
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
	game.cachePlayer(userId)

	game.sendPacketToPlayers(packets.NewServerGamePlayerWinCount(userId, playerWins.Wins))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetReferee Sets the referee for the game. Set userId to -1 to clear.
func (game *Game) SetReferee(requester *sessions.User, userId int) {
	if game.Data.InProgress {
		return
	}

	if !game.isUserHost(requester) {
		return
	}

	oldReferee := game.Data.RefereeId
	game.Data.RefereeId = userId

	game.spectators = utils.Filter(game.spectators, func(x int) bool { return x != oldReferee })

	if game.Data.RefereeId != -1 && !utils.Includes(game.spectators, userId) {
		game.spectators = append(game.spectators, userId)
	}

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
		game.cachePlayerScore(userId, score)
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
	game.validateAndCacheSettings()

	game.sendBotMessage(fmt.Sprintf("Tournament mode has been %v.", utils.BoolToEnabledString(game.Data.IsTournamentMode)))
	game.sendPacketToPlayers(packets.NewServerGameTournamentMode(game.Data.IsTournamentMode))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetClientProvidedDifficultyRatings Handles when the client provides new difficulty ratings for us to use
func (game *Game) SetClientProvidedDifficultyRatings(md5 string, alternativeMd5 string, difficulties []float64) {
	if len(difficulties) != countDifficultyRatings || len(game.Data.MapDifficultyRatingAll) == countDifficultyRatings {
		return
	}

	// Hash mismatch
	if md5 != game.Data.MapMD5 && alternativeMd5 != game.Data.MapMD5Alternative {
		return
	}

	game.Data.MapDifficultyRatingAll = difficulties
	game.validateAndCacheSettings()
	game.sendPacketToPlayers(packets.NewServerGameNeedDifficultyRatings(game.Data.MapMD5, game.Data.MapMD5Alternative, false))

	sendLobbyUsersGameInfoPacket(game, true)
}

// SetDonatorMapsetShared Sets whether an unsubmitted map is shared by a donator
func (game *Game) SetDonatorMapsetShared(isShared bool, sendToLobby bool) {
	game.Data.IsMapsetShared = isShared
	game.sendPacketToPlayers(packets.NewServerGameMapsetShared(isShared))

	if sendToLobby {
		sendLobbyUsersGameInfoPacket(game, true)
	}
}

// SetAutoHost Sets whether auto host is enabled for the game
func (game *Game) SetAutoHost(requester *sessions.User, enabled bool) {
	if !game.isUserHost(requester) {
		return
	}

	game.Data.IsAutoHost = enabled

	game.sendPacketToPlayers(packets.NewServerGameAutoHost(game.Data.IsAutoHost))
	sendLobbyUsersGameInfoPacket(game, true)

	if game.Data.IsAutoHost {
		game.sendBotMessage("Auto Host has been enabled. Use the following commands to further customize your game:\n" +
			"- `!mp mindiff (number)` - Changes the minimum difficulty that will be selected.\n" +
			"- `!mp maxdiff (number)` - Changes the maximum difficulty that will be selected.\n" +
			"- `!mp allowmode (4k/7k)` - Allows a game mode to be selected.\n" +
			"- `!mp disallowmode (4k/7k)` - Disallows a game mode to be selected.\n" +
			"- `!mp randmap` - Selects a new random map.")
		return
	}

	game.sendBotMessage(fmt.Sprintf("Auto Host has been disabled."))
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

// Handles disbandment of the multiplayer game
func (game *Game) disband() {
	game.EndGame()
	game.isDisbanded = true

	// Tournament mode games are kept around and deleted manually
	if game.Data.IsTournamentMode {
		return
	}

	game.deleteCachedMatchSettings()
	chat.RemoveMultiplayerChannel(game.Data.GameId)
	RemoveGameFromLobby(game)
}

// Returns if the user is host of the game or has permission.
func (game *Game) isUserHost(user *sessions.User) bool {
	if user == nil {
		return true
	}

	if user.Info.Id != game.Data.HostId && user.Info.Id != game.CreatorId && !common.HasUserGroup(user.Info.UserGroups, common.UserGroupDeveloper) {
		return false
	}

	return true
}

// Returns if a user is inside the game
func (game *Game) isUserInGame(user *sessions.User) bool {
	if user == nil {
		return false
	}

	return utils.Includes(game.Data.PlayerIds, user.Info.Id)
}

// Clears all players that are ready.
func (game *Game) clearReadyPlayers(sendToLobby bool) {
	for _, id := range game.Data.PlayersReady {
		game.cachePlayer(id)
		game.sendPacketToPlayers(packets.NewServerGamePlayerNotReady(id))
	}

	game.Data.PlayersReady = []int{}

	if sendToLobby {
		sendLobbyUsersGameInfoPacket(game, true)
	}
}

// Makes  the spectators start spectating playersInMatch.
func (game *Game) initializeSpectators() {
	for _, userId := range game.spectators {
		user := sessions.GetUserById(userId)
		game.initializeSpectator(user)
	}
}

func (game *Game) initializeSpectator(user *sessions.User) {
	if user == nil {
		return
	}

	user.StopSpectatingAll()

	if len(game.playersInMatch) != 2 && !common.HasUserGroup(user.Info.UserGroups, common.UserGroupDeveloper) {
		sessions.SendPacketToUser(packets.NewServerNotificationInfo("You can only spectate matches with two players."), user)
		return
	}

	for _, playerId := range game.playersInMatch {
		if player := sessions.GetUserById(playerId); player != nil {
			player.AddSpectator(user)
		}
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

// Returns the id of the most recent match winner.
func (game *Game) checkPlayerWinResult(userId int) (WinResult, error) {
	if _, ok := game.playerScores[userId]; !ok {
		return -1, errors.New("player score does not exist")
	}

	for scoreUserId, score := range game.playerScores {
		if scoreUserId == userId {
			continue
		}

		if game.playerScores[userId].PerformanceRating < score.PerformanceRating {
			return WinResultLost, nil
		}
	}

	return WinResultWon, nil
}

// Updates the win count for each player
func (game *Game) updatePlayerWinCount() {
	for userId := range game.playerScores {
		winResult, err := game.checkPlayerWinResult(userId)

		if err != nil {
			continue
		}

		if winResult != WinResultWon {
			continue
		}

		playerWins, err := utils.Find(game.Data.PlayerWins, func(x *objects.MultiplayerGamePlayerWins) bool { return x.Id == userId })

		if err != nil {
			continue
		}

		game.SetPlayerWinCount(userId, playerWins.Wins+1)
	}
}

// Inserts the current match into the database.
func (game *Game) insertMatchIntoDatabase() {
	if len(game.playerScores) == 0 {
		return
	}

	match := db.MultiplayerMatch{
		GameId:          game.Data.Id,
		TimePlayed:      time.Now().UnixMilli(),
		MapMd5:          game.Data.MapMD5,
		MapName:         game.Data.MapName,
		HostId:          game.Data.HostId,
		Ruleset:         game.Data.Ruleset,
		GameMode:        game.Data.MapGameMode,
		GlobalModifiers: game.Data.GlobalModifiers,
		FreeMod:         game.Data.FreeModType,
	}

	err := match.InsertIntoDatabase()

	if err != nil {
		log.Printf("Failed to insert match from game #%v into database - %v\n", game.Data.Id, err)
		return
	}

	for userId, score := range game.playerScores {
		winResult, _ := game.checkPlayerWinResult(userId)

		dbScore := db.MultiplayerMatchScore{
			UserId:            userId,
			MatchId:           match.Id,
			Mods:              score.Modifiers,
			PerformanceRating: score.PerformanceRating,
			Accuracy:          score.Accuracy,
			MaxCombo:          score.MaxCombo,
			CountMarv:         score.Judgements[common.JudgementMarv],
			CountPerf:         score.Judgements[common.JudgementPerf],
			CountGreat:        score.Judgements[common.JudgementGreat],
			CountGood:         score.Judgements[common.JudgementGood],
			CountOkay:         score.Judgements[common.JudgementOkay],
			CountMiss:         score.Judgements[common.JudgementMiss],
			Won:               int(winResult),
		}

		err := dbScore.InsertIntoDatabase()

		if err != nil {
			log.Printf("Failed to insert #%v's match score from game #%v into database - %v\n", userId, game.Data.Id, err)
			return
		}
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

func (game *Game) isPlayerSpectatorOrReferee(userId int) bool {
	return game.Data.RefereeId == userId || slices.Contains(game.spectators, userId)
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

// Selects a random map from the database according to difficulty filters
func (game *Game) selectAutohostMap() {
	song, err := db.GetRandomSongMap(game.Data.FilterMinDifficultyRating, game.Data.FilterMaxDifficultyRating, game.Data.FilterAllowedGameModes)

	if err != nil {
		log.Printf("error selecting random map in multiplayer - %v\n", err)
		return
	}

	game.changeMapFromDbSong(song)
}

func (game *Game) changeMapFromDbSong(song *db.SongMap) {
	game.ChangeMap(nil, &packets.ClientChangeGameMap{
		MD5:                 song.Md5.String,
		AlternativeMD5:      song.AlternativeMd5.String,
		MapId:               song.Id,
		MapsetId:            song.MapsetId,
		Name:                fmt.Sprintf("%v - %v [%v]", song.Artist.String, song.Title.String, song.DifficultyName.String),
		Mode:                song.GameMode,
		DifficultyRating:    song.DifficultyRating,
		DifficultyRatingAll: []float64{},
	})
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

	for _, id := range game.spectators {
		// Referee will have already gotten the packet above.
		if id == game.Data.RefereeId && slices.Contains(game.Data.PlayerIds, id) {
			continue
		}

		user := sessions.GetUserById(id)

		if user == nil {
			continue
		}

		sessions.SendPacketToUser(packet, user)
	}
}

// Sends a message to the multiplayer chat from the bot
func (game *Game) sendBotMessage(message string) {
	chat.SendMessage(chat.Bot, game.chatChannel.Name, message)
}

// validateAndCacheSettings Checks the multiplayer settings to see if they are in an acceptable range
func (game *Game) validateAndCacheSettings() {
	data := game.Data

	data.Name = utils.TruncateString(data.Name, 50)

	if censored := utils.CensorString(data.Name); censored != "" {
		data.Name = censored
	}

	data.HasPassword = game.Password != ""
	data.MaxPlayers = utils.Clamp(data.MaxPlayers, 2, 16)
	data.Ruleset = objects.MultiplayerGameRulesetFreeForAll
	data.FreeModType = utils.Clamp(data.FreeModType, objects.MultiplayerGameFreeModNone, objects.MultiplayerGameFreeModRegular|objects.MultiplayerGameFreeModRate)

	data.MapMD5 = utils.TruncateString(data.MapMD5, 64)
	data.MapMD5Alternative = utils.TruncateString(data.MapMD5Alternative, 64)
	data.MapName = utils.TruncateString(data.MapName, 250)
	data.MapGameMode = utils.Clamp(data.MapGameMode, common.ModeKeys4, common.ModeKeys7)

	data.FilterMinDifficultyRating = utils.Clamp(data.FilterMinDifficultyRating, 0, 100)
	data.FilterMaxDifficultyRating = utils.Clamp(data.FilterMaxDifficultyRating, 0, 100)
	data.FilterMaxSongLength = utils.Clamp(data.FilterMaxSongLength, 0, math.MaxInt32)
	data.FilterMinLongNotePercent = utils.Clamp(data.FilterMinLongNotePercent, 0, 100)
	data.FilterMaxLongNotePercent = utils.Clamp(data.FilterMaxLongNotePercent, 0, 100)
	data.FilterMinAudioRate = utils.Clamp(data.FilterMinAudioRate, 0.5, 2.0)

	// There is a maximum of 31 rates allowed in the game. So if we don't have all of them, then just clear it.
	if len(data.MapDifficultyRatingAll) != countDifficultyRatings {
		data.NeedsDifficultyRatings = true
		data.MapDifficultyRatingAll = []float64{}
		game.sendPacketToPlayers(packets.NewServerGameNeedDifficultyRatings(data.MapMD5, data.MapMD5Alternative, data.NeedsDifficultyRatings))
	} else {
		data.NeedsDifficultyRatings = false
	}

	if len(data.FilterAllowedGameModes) == 0 {
		data.FilterAllowedGameModes = []common.Mode{common.ModeKeys4, common.ModeKeys7}
	}

	game.cacheMatchSettings()
}

// Removes inactive players from the game
func (game *Game) removeInactivePlayers() {
	go func() {
		for !game.isDisbanded {
			game.RunLocked(func() {
				playerIds := make([]int, len(game.Data.PlayerIds))
				copy(playerIds, game.Data.PlayerIds)

				for _, playerId := range playerIds {
					user := sessions.GetUserById(playerId)

					if user != nil {
						continue
					}

					game.RemovePlayer(playerId)
					log.Printf("Removing %v from game: %v (%v)\n", playerId, game.Data.Name, game.Data.Id)
				}
			})

			time.Sleep(time.Second * 5)
		}
	}()
}
