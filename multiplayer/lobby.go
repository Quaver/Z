package multiplayer

import (
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
}

// RemoveUserFromLobby Removes a user from the multiplayer lobby
func RemoveUserFromLobby(user *sessions.User) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()

	delete(lobby.users, user.Info.Id)
}

// AddGameToLobby Adds a game to the multiplayer lobby list
// TODO: Send game info to lobby players
// TODO: Place user in game
func AddGameToLobby(game *Game) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()

	lobby.games[game.Data.Id] = game
	log.Printf("Multiplayer Game `%v (#%v)` was created.\n", game.Data.Name, game.Data.Id)
}

// RemoveGameFromLobby Removes a game from the multiplayer lobby list
// TODO: Remove all players from the game.
func RemoveGameFromLobby(game *Game) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()

	delete(lobby.games, game.Data.Id)
	log.Printf("Multiplayer game `%v (%v)` was disbanded.\n", game.Data.Name, game.Data.Id)
}
