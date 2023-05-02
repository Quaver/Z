package sessions

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/utils"
	"fmt"
	"net"
	"sync"
	"time"
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

	// The last time the user was pinged
	LastPingTimestamp int64

	// The last time the user sent a successful pong
	LastPongTimestamp int64

	// The last detected processes that were discovered on the user
	LastDetectedProcesses []int
}

type PacketUser struct {
	Id          int               `json:"id"`
	SteamId     string            `json:"sid"`
	Username    string            `json:"u"`
	UserGroups  common.UserGroups `json:"ug"`
	MuteEndTime int64             `json:"m"`
	Country     string            `json:"c"`
}

// NewUser Creates a new user session struct object
func NewUser(conn net.Conn, user *db.User) *User {
	return &User{
		Conn:              conn,
		Token:             utils.GenerateRandomString(64),
		Info:              user,
		Mutex:             &sync.Mutex{},
		Stats:             map[common.Mode]*db.UserStats{},
		LastPingTimestamp: time.Now().UnixMilli(),
		LastPongTimestamp: time.Now().UnixMilli(),
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

// SerializeForPacket Serializes the user to be used in a packet
func (u *User) SerializeForPacket() *PacketUser {
	return &PacketUser{
		Id:          u.Info.Id,
		SteamId:     u.Info.SteamId,
		Username:    u.Info.Username,
		UserGroups:  u.Info.UserGroups,
		MuteEndTime: u.Info.MuteEndTime,
		Country:     u.Info.Country,
	}
}

// Retrieves the Redis key for the user's session
func (u *User) getRedisSessionKey() string {
	return fmt.Sprintf("quaver:server:session:%v", u.Token)
}
