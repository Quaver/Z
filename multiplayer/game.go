package multiplayer

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/objects"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
	"math"
	"sync"
	"time"
)

type Game struct {
	Data           *objects.MultiplayerGame
	Password       string      // The password for the game. This is different from Data.CreationPassword, as it is hidden from users.
	CreatorId      int         // The id of the user who created the game
	mutex          *sync.Mutex // Locks down the game to prevent race conditions
	countdownTimer *time.Timer // Counts down before starting the game
}

const (
	maxPlayerCount int = 16 // The maximum amount of players allowed in a game
)

// NewGame Creates a new multiplayer game from a game
func NewGame(gameData *objects.MultiplayerGame, creatorId int) (*Game, error) {
	game := Game{
		Data:      gameData,
		CreatorId: creatorId,
		mutex:     &sync.Mutex{},
		Password:  gameData.CreationPassword,
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

	return &game, nil
}

// AddPlayer Adds a user to the multiplayer game
func (game *Game) AddPlayer(userId int, password string) {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	user := sessions.GetUserById(userId)

	if user == nil {
		return
	}

	currentGame := GetGameById(user.GetMultiplayerGameId())

	if currentGame != nil && currentGame.Data.Id != game.Data.Id {
		currentGame.RemovePlayer(user.Info.Id)
	}

	if len(game.Data.PlayerIds) >= maxPlayerCount {
		sessions.SendPacketToUser(packets.NewServerJoinGameFailed(packets.JoinGameErrorFull), user)
		return
	}

	if game.Data.HasPassword && game.Password != password && !common.HasUserGroup(user.Info.UserGroups, common.UserGroupSwan) {
		sessions.SendPacketToUser(packets.NewServerJoinGameFailed(packets.JoinGameErrorPassword), user)
		return
	}

	game.Data.PlayerIds = append(game.Data.PlayerIds, user.Info.Id)
	game.Data.PlayerModifiers = append(game.Data.PlayerModifiers, objects.MultiplayerGamePlayerMods{Id: user.Info.Id})
	game.Data.PlayerWins = append(game.Data.PlayerWins, objects.MultiplayerGamePlayerWins{Id: user.Info.Id})

	user.SetMultiplayerGameId(game.Data.Id)
	game.sendPacketToPlayers(packets.NewServerUserJoinedGame(user.Info.Id))
	sessions.SendPacketToUser(packets.NewServerJoinGame(game.Data.GameId), user)

	if len(game.Data.PlayerIds) == 1 {
		game.SetHost(user.Info.Id, false)
	}

	RemoveUserFromLobby(user)
	sendLobbyUsersGameInfoPacket(game, true)
}

// RemovePlayer Removes a player from the multiplayer game and disbands the game if necessary
func (game *Game) RemovePlayer(userId int) {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	user := sessions.GetUserById(userId)

	if user == nil {
		return
	}

	user.SetMultiplayerGameId(0)

	game.Data.PlayerIds = utils.Filter(game.Data.PlayerIds, func(x int) bool { return x != user.Info.Id })
	game.Data.PlayerModifiers = utils.Filter(game.Data.PlayerModifiers, func(x objects.MultiplayerGamePlayerMods) bool { return x.Id != user.Info.Id })
	game.Data.PlayerWins = utils.Filter(game.Data.PlayerWins, func(x objects.MultiplayerGamePlayerWins) bool { return x.Id != user.Info.Id })

	// Disband game since there are no more players left
	if len(game.Data.PlayerIds) == 0 {
		//RemoveGameFromLobby(game)
		return
	}

	game.SetHost(game.Data.PlayerIds[0], false)
	game.sendPacketToPlayers(packets.NewServerUserLeftGame(user.Info.Id))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetHost Sets the host of the game
func (game *Game) SetHost(userId int, lock bool) {
	if lock {
		game.mutex.Lock()
		defer game.mutex.Unlock()
	}

	user := sessions.GetUserById(userId)

	if user == nil {
		return
	}

	if user.GetMultiplayerGameId() != game.Data.Id {
		return
	}

	game.Data.HostId = userId
	game.SetHostSelectingMap(nil, false, false, false)

	game.sendPacketToPlayers(packets.NewServerGameChangeHost(game.Data.HostId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// ChangeMap Changes the multiplayer map. Non-nil requester checks if they are the host
func (game *Game) ChangeMap(requester *sessions.User, packet *packets.ClientChangeGameMap) {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	if game.Data.InProgress {
		return
	}

	if requester != nil && requester.Info.Id != game.Data.HostId {
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
	game.validateSettings()
	game.clearReadyPlayers(false)

	game.sendPacketToPlayers(packets.NewServerGameMapChanged(packet))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerDoesntHaveMap Sets that a player does not have the map downloaded
func (game *Game) SetPlayerDoesntHaveMap(userId int) {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	game.Data.PlayersWithoutMap = append(game.Data.PlayersWithoutMap, userId)

	game.sendPacketToPlayers(packets.NewServerGamePlayerNoMap(userId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerHasMap Sets that a player now has the currently played map
func (game *Game) SetPlayerHasMap(userId int) {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	game.Data.PlayersWithoutMap = utils.Filter(game.Data.PlayersWithoutMap, func(x int) bool {
		return x != userId
	})

	game.sendPacketToPlayers(packets.NewServerGamePlayerHasMap(userId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerReady Sets that a player is currently ready to play
func (game *Game) SetPlayerReady(userId int) {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	if !utils.Includes(game.Data.PlayersReady, userId) {
		game.Data.PlayersReady = append(game.Data.PlayersReady, userId)
	}

	game.sendPacketToPlayers(packets.NewServerGamePlayerReady(userId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// SetPlayerNotReady Sets that a player is not ready to play
func (game *Game) SetPlayerNotReady(userId int) {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	game.Data.PlayersReady = utils.Filter(game.Data.PlayersReady, func(i int) bool {
		return i != userId
	})

	game.sendPacketToPlayers(packets.NewServerGamePlayerNotReady(userId))
	sendLobbyUsersGameInfoPacket(game, true)
}

// StartCountdown Starts the 5-second multiplayer countdown
func (game *Game) StartCountdown(requester *sessions.User) {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	if game.Data.InProgress {
		return
	}

	if requester != nil && requester.Info.Id != game.Data.HostId {
		return
	}

	game.countdownTimer = time.AfterFunc(5*time.Second, func() {
		game.StartGame()
	})

	game.sendPacketToPlayers(packets.NewServerGameStartCountdown())
	sendLobbyUsersGameInfoPacket(game, true)
}

// StopCountdown Stops the multiplayer countdown if one is live
func (game *Game) StopCountdown(requester *sessions.User) {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	if game.Data.InProgress {
		return
	}

	if requester != nil && requester.Info.Id != game.Data.HostId {
		return
	}

	game.clearCountdown()
	sendLobbyUsersGameInfoPacket(game, true)
}

// StartGame Starts the multiplayer game
func (game *Game) StartGame() {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	if game.Data.InProgress {
		return
	}

	game.clearCountdown()
	game.clearReadyPlayers(false)

	sendLobbyUsersGameInfoPacket(game, true)
}

// ChangeName Changes the name of the multiplayer game
func (game *Game) ChangeName(requester *sessions.User, name string) {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	if game.Data.InProgress {
		return
	}

	if requester != nil && requester.Info.Id != game.Data.HostId {
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
func (game *Game) SetHostSelectingMap(requester *sessions.User, isSelecting bool, sendToLobby bool, lock bool) {
	if lock {
		game.mutex.Lock()
		defer game.mutex.Unlock()
	}

	if game.Data.InProgress {
		return
	}

	if requester != nil && requester.Info.Id != game.Data.HostId {
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
	game.mutex.Lock()
	defer game.mutex.Unlock()

	if requester != nil && requester.Info.Id != game.Data.HostId {
		return
	}

	game.Password = password
	game.validateSettings()

	sendLobbyUsersGameInfoPacket(game, true)
}

// Clears all players that are ready. This is to be used in an already mutex-locked context.
func (game *Game) clearReadyPlayers(sendToLobby bool) {
	for _, id := range game.Data.PlayersReady {
		game.sendPacketToPlayers(packets.NewServerGamePlayerNotReady(id))
	}

	game.Data.PlayersReady = []int{}

	if sendToLobby {
		sendLobbyUsersGameInfoPacket(game, true)
	}
}

// Clears and stops the countdown timer This is to be used in an already mutex-locked context.
func (game *Game) clearCountdown() {
	if game.countdownTimer != nil {
		game.countdownTimer.Stop()
		game.countdownTimer = nil
	}

	game.sendPacketToPlayers(packets.NewServerGameStopCountdown())
}

// Sends a packet to all players in the game. This is to be used in an already mutex-locked context.
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
	data.Ruleset = utils.Clamp(data.Ruleset, objects.MultiplayerGameRulesetFreeForAll, objects.MultiplayerGameRulesetTeam)
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

	// There is a maximum of 21 rates allowed in the game. So if we don't have all of them, then just clear it.
	if len(data.MapDifficultyRatingAll) < 21 {
		data.MapDifficultyRatingAll = []float64{}
	}
}
