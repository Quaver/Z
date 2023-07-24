package sessions

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/objects"
	"example.com/Quaver/Z/utils"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type User struct {
	// The connection for the user
	Conn net.Conn

	// The token used to identify the user for requests.
	token string

	// All user table information from the database
	Info *db.User

	// mutex for all operations regarding changes in the user
	mutex *sync.Mutex

	// Player statistics from the database
	stats map[common.Mode]*db.UserStats

	// The last time the user was pinged
	lastPingTimestamp int64

	// The last time the user sent a successful pong
	lastPongTimestamp int64

	// The last detected processes that were discovered on the user
	lastDetectedProcesses []string

	// The current client status of the user
	status *objects.ClientStatus

	// A count of the amount of messages the user has spammed in the past x amount of time. Used for muting purposes.
	spammedChatMessages int

	// The last time the user's spammed messages were checked
	spammedChatLastTimeCleared int64

	// The id of the multiplayer game if the user is inside of one
	multiplayerGameId int
}

// NewUser Creates a new user session struct object
func NewUser(conn net.Conn, user *db.User) *User {
	return &User{
		Conn:              conn,
		token:             utils.GenerateRandomString(64),
		Info:              user,
		mutex:             &sync.Mutex{},
		stats:             map[common.Mode]*db.UserStats{},
		lastPingTimestamp: time.Now().UnixMilli(),
		lastPongTimestamp: time.Now().UnixMilli(),
		status: &objects.ClientStatus{
			Status:    0,
			MapId:     -1,
			MapMd5:    "",
			GameMode:  common.ModeKeys4,
			Content:   "",
			Modifiers: 0,
		},
	}
}

// GetToken Returns the user token
func (u *User) GetToken() string {
	return u.token
}

// GetStats Retrieves the stats for the user
func (u *User) GetStats() map[common.Mode]*db.UserStats {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.stats
}

func (u *User) GetStatsSlice() []*db.PacketUserStats {
	statSlice := make([]*db.PacketUserStats, 0)

	for _, value := range u.GetStats() {
		statSlice = append(statSlice, value.SerializeForPacket())
	}

	return statSlice
}

// SetStats Updates the statistics for the user
func (u *User) SetStats() error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	for i := 1; i < int(common.ModeEnumMaxValue); i++ {
		mode := common.Mode(i)
		stats, err := db.GetUserStats(u.Info.Id, u.Info.Country, mode)

		if err != nil {
			return err
		}

		u.stats[mode] = stats
	}

	return nil
}

// GetLastPingTimestamp Retrieves the last ping timestamp
func (u *User) GetLastPingTimestamp() int64 {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.lastPingTimestamp
}

// SetLastPingTimestamp Sets the last ping timestamp to the current time
func (u *User) SetLastPingTimestamp() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.lastPingTimestamp = time.Now().UnixMilli()
}

// GetLastPongTimestamp Retrieves the last pong timestamp
func (u *User) GetLastPongTimestamp() int64 {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.lastPongTimestamp
}

// SetLastPongTimestamp Sets the last pong timestamp to the current time
func (u *User) SetLastPongTimestamp() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.lastPongTimestamp = time.Now().UnixMilli()
}

// GetLastDetectedProcesses Gets the last detected processes for the user
func (u *User) GetLastDetectedProcesses() []string {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.lastDetectedProcesses
}

// SetLastDetectedProcesses Sets the last detected processes for the user
func (u *User) SetLastDetectedProcesses(processes []string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.lastDetectedProcesses = processes
}

// GetClientStatus Gets the current user client status
func (u *User) GetClientStatus() *objects.ClientStatus {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.status
}

// SetClientStatus Sets the current user client status
func (u *User) SetClientStatus(status *objects.ClientStatus) {
	u.mutex.Lock()
	u.status = status
	u.mutex.Unlock()

	err := addUserClientStatusToRedis(u)

	if err != nil {
		log.Println(err)
	}
}

// GetSpammedMessagesCount Gets the amount of messages the user has spammed
func (u *User) GetSpammedMessagesCount() int {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.spammedChatMessages
}

// IncrementSpammedMessagesCount Increments the amount of spammed messages by 1.
func (u *User) IncrementSpammedMessagesCount() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.spammedChatMessages++
}

// ResetSpammedMessagesCount Resets the amount of spammed messages back to zero
func (u *User) ResetSpammedMessagesCount() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.spammedChatMessages = 0
}

// GetSpammedChatLastTimeCleared Gets the last time the user's chat spam rate was cleared
func (u *User) GetSpammedChatLastTimeCleared() int64 {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.spammedChatLastTimeCleared
}

// SetSpammedChatLastTimeCleared Sets the time the user's chat spam rate was cleared
func (u *User) SetSpammedChatLastTimeCleared(time int64) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.spammedChatLastTimeCleared = time
}

// GetMultiplayerGameId Gets the id of the multiplayer game the user is currently inside of (if any)
func (u *User) GetMultiplayerGameId() int {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.multiplayerGameId
}

// SetMultiplayerGameId Sets the id of the multiplayer game if the user is inside of one
func (u *User) SetMultiplayerGameId(id int) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.multiplayerGameId = id
}

// IsMuted Returns if the user is muted
func (u *User) IsMuted() bool {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.Info.MuteEndTime > time.Now().UnixMilli()
}

// MuteUser Mutes a user for a specified duration
func (u *User) MuteUser(duration time.Duration) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	endTime := time.Now().UnixMilli() + duration.Milliseconds()

	err := db.MuteUser(u.Info.Id, endTime)

	if err != nil {
		log.Printf("Failed to update user mute time: %v\n", err)
		return err
	}

	u.Info.MuteEndTime = endTime
	return nil
}

// SerializeForPacket Serializes the user to be used in a packet
func (u *User) SerializeForPacket() *objects.PacketUser {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return &objects.PacketUser{
		Id:          u.Info.Id,
		SteamId:     u.Info.SteamId,
		Username:    u.Info.Username,
		UserGroups:  u.Info.UserGroups,
		MuteEndTime: u.Info.MuteEndTime,
		Country:     u.Info.Country,
	}
}

// Returns the Redis key for the user's session
func (u *User) getRedisSessionKey() string {
	return fmt.Sprintf("quaver:server:session:%v", u.token)
}

// Returns the Redis key for the user's client
func (u *User) getRedisClientStatusKey() string {
	return fmt.Sprintf("quaver:server:user_status:%v", u.Info.Id)
}
