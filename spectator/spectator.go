package spectator

import (
	"example.com/Quaver/Z/sessions"
	"sync"
)

var (
	mutex        = &sync.Mutex{}
	userIdToUser = map[int]*User{}
)

// GetUser Returns a user by their session.User
func GetUser(user *sessions.User) *User {
	return GetUserById(user.Info.Id)
}

// GetUserById Returns a spectator user by their id
func GetUserById(id int) *User {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := userIdToUser[id]; ok {
		return userIdToUser[id]
	}

	return nil
}

// InitializeUser Initializes a user to able to be spectated
func InitializeUser(user *sessions.User) {
	if GetUser(user) != nil {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	userIdToUser[user.Info.Id] = &User{User: user}
}

// UninitializeUser Uninitializes a user to be spectated
func UninitializeUser(user *sessions.User) {
	GetUserById(user.Info.Id).StopSpectatingAll()

	mutex.Lock()
	defer mutex.Unlock()
	delete(userIdToUser, user.Info.Id)
}
