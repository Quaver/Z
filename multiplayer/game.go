package multiplayer

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/objects"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
	"log"
	"math"
	"sync"
)

type Game struct {
	Data      *objects.MultiplayerGame
	Password  string // The password for the game. This is different from Data.CreationPassword, as it is hidden from users.
	CreatorId int    // The id of the user who created the game
	mutex     *sync.Mutex
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

	// Remove from player ids
	// Remove from player modifiers
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
		log.Printf("[MP #%v] Tried to give host to user `%v (%v)`, but they are not in the game\n", game.Data.Id, user.Info.Username, user.Info.Id)
		return
	}

	game.Data.HostId = userId

	game.sendPacketToPlayers(packets.NewServerGameChangeHost(game.Data.HostId))
	sendLobbyUsersGameInfoPacket(game, true)
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
	data.Ruleset = utils.Clamp(data.Ruleset, objects.MultiplayerGameRulesetFreeForAll, objects.MultiplayerGameRulesetTeam)
	data.FreeModType = utils.Clamp(data.FreeModType, objects.MultiplayerGameFreeModNone, objects.MultiplayerGameFreeModRegular|objects.MultiplayerGameFreeModRate)

	data.MapMD5 = utils.TruncateString(data.MapMD5, 32)
	data.MapMD5Alternative = utils.TruncateString(data.MapMD5Alternative, 32)
	data.MapName = utils.TruncateString(data.MapName, 250)
	data.MapGameMode = utils.Clamp(data.MapGameMode, common.ModeKeys4, common.ModeKeys7)

	data.FilterMinDifficultyRating = utils.Clamp(data.FilterMinDifficultyRating, 0, math.MaxInt32)
	data.FilterMaxDifficultyRating = utils.Clamp(data.FilterMaxDifficultyRating, 0, math.MaxInt32)
	data.FilterMaxSongLength = utils.Clamp(data.FilterMaxSongLength, 0, math.MaxInt32)
	data.FilterMinLongNotePercent = utils.Clamp(data.FilterMinLongNotePercent, 0, 100)
	data.FilterMaxLongNotePercent = utils.Clamp(data.FilterMaxLongNotePercent, 0, 100)
	data.FilterMinAudioRate = utils.Clamp(data.FilterMinAudioRate, 0.5, 2.0)

	// There is a maximum of 21 rates allowed in the game. So if we don't have all of them, then just clear it.
	if len(data.MapAllDifficultyRatings) < 21 {
		data.MapAllDifficultyRatings = []float64{}
	}
}
