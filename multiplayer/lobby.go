package multiplayer

import (
	"example.com/Quaver/Z/sessions"
	"log"
	"sync"
)

type multiplayerLobby struct {
	users map[int]*sessions.User
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
		mutex: &sync.Mutex{},
	}
}

// AddUserToLobby Adds a user to the multiplayer lobby
func AddUserToLobby(user *sessions.User) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()

	lobby.users[user.Info.Id] = user
	log.Printf("[%v - #%v] joined the multiplayer lobby\n", user.Info.Username, user.Info.Id)
}

// RemoveUserFromLobby Removes a user from the multiplayer lobby
func RemoveUserFromLobby(user *sessions.User) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()

	delete(lobby.users, user.Info.Id)
	log.Printf("[%v - #%v] left the multiplayer lobby\n", user.Info.Username, user.Info.Id)
}
