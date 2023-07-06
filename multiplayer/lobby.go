package multiplayer

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"log"
	"sync"
)

type multiplayerLobby struct {
	users map[int]*sessions.User
	games map[int]*Game
	mutex *sync.Mutex
}

var lobby *multiplayerLobby

// InitializeLobby Initializes the multiplayer lobby / games
func InitializeLobby() {
	if lobby != nil {
		return
	}

	lobby = &multiplayerLobby{
		users: map[int]*sessions.User{},
		games: map[int]*Game{},
		mutex: &sync.Mutex{},
	}
}

// AddUserToLobby Adds a user to the multiplayer lobby
func AddUserToLobby(user *sessions.User) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()

	lobby.users[user.Info.Id] = user

	for _, game := range lobby.games {
		sendLobbyUsersGameInfoPacket(game, false)
	}
}

// RemoveUserFromLobby Removes a user from the multiplayer lobby
func RemoveUserFromLobby(user *sessions.User) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()

	delete(lobby.users, user.Info.Id)
}

// AddGameToLobby Adds a game to the multiplayer lobby list
func AddGameToLobby(game *Game) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()

	lobby.games[game.Data.Id] = game
	sendLobbyUsersGameInfoPacket(game, false)

	log.Printf("Multiplayer Game `%v (#%v)` was created.\n", game.Data.Name, game.Data.Id)
}

// RemoveGameFromLobby Removes a game from the multiplayer lobby list
func RemoveGameFromLobby(game *Game) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()

	for _, user := range lobby.users {
		sessions.SendPacketToUser(packets.NewServerGameDisbanded(game.Data.GameId), user)
	}

	delete(lobby.games, game.Data.Id)
	log.Printf("Multiplayer game `%v (%v)` was disbanded.\n", game.Data.Name, game.Data.Id)
}

// GetGameById Retrieves a multiplayer game by its id
func GetGameById(id int) *Game {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()

	return lobby.games[id]
}

// GetGameByIdString Retrieves a game by its stringified id
func GetGameByIdString(id string) *Game {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()

	for _, game := range lobby.games {
		if game.Data.GameId == id {
			return game
		}
	}

	return nil
}

// SendLobbyUsersGameInfoPacket Sends all the users in the lobby a packet with game information
// Be careful of deadlocks when calling this. Make sure not to call the mutex twice.
func sendLobbyUsersGameInfoPacket(game *Game, lock bool) {
	if lock {
		lobby.mutex.Lock()
		defer lobby.mutex.Unlock()
	}

	packet := packets.NewServerMultiplayerGameInfo(game.Data)

	for _, user := range lobby.users {
		sessions.SendPacketToUser(packet, user)
	}
}
