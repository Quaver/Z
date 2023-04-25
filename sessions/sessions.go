package sessions

import (
	"example.com/Quaver/Z/db"
	"net"
	"sync"
)

var (
	// Mutex used for thread-safe access to users
	userMutex = &sync.Mutex{}

	// A map to users with the key being their user id
	userIdToUser = map[int]*User{}

	// A map to users with the key being their username
	usernameToUser = map[string]*User{}

	// A map to users with the key being their connection
	connToUser = map[net.Conn]*User{}
)

// AddUser Adds a user session
func AddUser(user *User) error {
	addUserToMaps(user)

	err := UpdateRedisOnlineUserCount()

	if err != nil {

		return err
	}

	return nil
}

// RemoveUser Removes a user session
func RemoveUser(user *User) error {
	removeUserFromMaps(user)

	err := UpdateRedisOnlineUserCount()

	if err != nil {
		return err
	}

	return nil
}

// GetUserById Returns a user by their id
func GetUserById(id int) *User {
	userMutex.Lock()
	defer userMutex.Unlock()

	return userIdToUser[id]
}

// GetUserByUsername Returns a user by their username
func GetUserByUsername(username string) *User {
	userMutex.Lock()
	defer userMutex.Unlock()

	return usernameToUser[username]
}

// GetUserByConnection Returns a user by their connection to the server
func GetUserByConnection(conn net.Conn) *User {
	userMutex.Lock()
	defer userMutex.Unlock()

	return connToUser[conn]
}

// GetOnlineUserCount Returns the number of online users
func GetOnlineUserCount() int {
	userMutex.Lock()
	defer userMutex.Unlock()

	return len(userIdToUser)
}

// Adds a user to the maps that can be used to look them up
func addUserToMaps(user *User) {
	userMutex.Lock()
	defer userMutex.Unlock()

	userIdToUser[user.Info.Id] = user
	usernameToUser[user.Info.Username] = user
	connToUser[user.Conn] = user
}

// Removes a user from the maps that are used to look them up
func removeUserFromMaps(user *User) {
	userMutex.Lock()
	defer userMutex.Unlock()

	delete(userIdToUser, user.Info.Id)
	delete(usernameToUser, user.Info.Username)
	delete(connToUser, user.Conn)
}

// UpdateRedisOnlineUserCount Updates the online user count in Redis
func UpdateRedisOnlineUserCount() error {
	_, err := db.Redis.Set(db.RedisCtx, "quaver:server:online_users", GetOnlineUserCount(), 0).Result()

	if err != nil {
		return err
	}

	return nil
}
