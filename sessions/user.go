package sessions

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/utils"
	"net"
	"sync"
)

type User struct {
	// The connection for the user
	Conn net.Conn

	// The token used to identify the user for requests.
	Token string

	// All user table information from the database
	Info *db.User

	// Mutex for all operations regarding changes in the user
	Mutex *sync.Mutex

	// Player statistics from the database
	Stats map[common.Mode]*db.UserStats
}

// NewUser Creates a new user session struct object
func NewUser(conn net.Conn, user *db.User) *User {
	return &User{
		Conn:  conn,
		Token: utils.GenerateRandomString(64),
		Info:  user,
		Mutex: &sync.Mutex{},
		Stats: map[common.Mode]*db.UserStats{},
	}
}

// UpdateStats Updates the statistics for the user
func (u *User) UpdateStats() error {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	for i := 1; i < int(common.ModeEnumMaxValue); i++ {
		mode := common.Mode(i)
		stats, err := db.GetUserStats(u.Info.Id, u.Info.Country, mode)

		if err != nil {
			return err
		}

		u.Stats[mode] = stats
	}

	return nil
}
