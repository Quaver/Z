package sessions

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/objects"
	"example.com/Quaver/Z/packets"
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

	// If the user's session is fully closed and logged out.
	SessionClosed bool

	// The token used to identify the user for requests.
	token string

	// All user table information from the database
	Info *db.User

	// Mutex for all operations regarding changes in the user
	Mutex *sync.Mutex

	// Player statistics from the database
	stats map[common.Mode]*db.UserStats

	// The last time the user was pinged
	lastPingTimestamp int64

	// The last time the user sent a successful pong
	lastPongTimestamp int64

	lastTemporaryDisconnectionTimestamp int64

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

	// People who are currently watching this user
	spectators []*User

	// People who the user is currently watching
	spectating []*User

	// The replay frames for the user's current play session
	frames []*packets.ClientSpectatorReplayFrames
}

// NewUser Creates a new user session struct object
func NewUser(conn net.Conn, user *db.User) *User {
	sessionUser := User{
		Conn:                                conn,
		SessionClosed:                       false,
		token:                               utils.GenerateRandomString(64),
		Info:                                user,
		Mutex:                               &sync.Mutex{},
		stats:                               map[common.Mode]*db.UserStats{},
		lastPingTimestamp:                   time.Now().UnixMilli(),
		lastPongTimestamp:                   time.Now().UnixMilli(),
		lastTemporaryDisconnectionTimestamp: -1,
		status: &objects.ClientStatus{
			Status:    0,
			MapId:     -1,
			MapMd5:    "",
			GameMode:  common.ModeKeys4,
			Content:   "",
			Modifiers: 0,
		},
		spectators: []*User{},
		spectating: []*User{},
		frames:     []*packets.ClientSpectatorReplayFrames{},
	}

	return &sessionUser
}

// GetToken Returns the user token
func (u *User) GetToken() string {
	return u.token
}

// GetStats Retrieves the stats for the user
func (u *User) GetStats() map[common.Mode]*db.UserStats {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

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
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

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
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	return u.lastPingTimestamp
}

// SetLastPingTimestamp Sets the last ping timestamp to the current time
func (u *User) SetLastPingTimestamp() {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	u.lastPingTimestamp = time.Now().UnixMilli()
}

// GetLastPongTimestamp Retrieves the last pong timestamp
func (u *User) GetLastPongTimestamp() int64 {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	return u.lastPongTimestamp
}

// SetLastPongTimestamp Sets the last pong timestamp to the current time
func (u *User) SetLastPongTimestamp() {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	u.lastPongTimestamp = time.Now().UnixMilli()
}

func (u *User) GetLastTemporaryDisconnectionTimestamp() int64 {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	return u.lastTemporaryDisconnectionTimestamp
}

func (u *User) SetLastTemporaryDisconnectionTimestamp(lastTemporaryDisconnectionTimestamp int64) {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	u.lastTemporaryDisconnectionTimestamp = lastTemporaryDisconnectionTimestamp
}

// GetLastDetectedProcesses Gets the last detected processes for the user
func (u *User) GetLastDetectedProcesses() []string {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	return u.lastDetectedProcesses
}

// SetLastDetectedProcesses Sets the last detected processes for the user
func (u *User) SetLastDetectedProcesses(processes []string) {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	u.lastDetectedProcesses = processes
}

// GetClientStatus Gets the current user client status
func (u *User) GetClientStatus() *objects.ClientStatus {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	return u.status
}

// SetClientStatus Sets the current user client status
func (u *User) SetClientStatus(status *objects.ClientStatus) {
	u.Mutex.Lock()
	u.status = status
	u.Mutex.Unlock()

	err := addUserClientStatusToRedis(u)

	if err != nil {
		log.Println(err)
	}
}

// GetSpammedMessagesCount Gets the amount of messages the user has spammed
func (u *User) GetSpammedMessagesCount() int {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	return u.spammedChatMessages
}

// IncrementSpammedMessagesCount Increments the amount of spammed messages by 1.
func (u *User) IncrementSpammedMessagesCount() {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	u.spammedChatMessages++
}

// ResetSpammedMessagesCount Resets the amount of spammed messages back to zero
func (u *User) ResetSpammedMessagesCount() {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	u.spammedChatMessages = 0
}

// GetSpammedChatLastTimeCleared Gets the last time the user's chat spam rate was cleared
func (u *User) GetSpammedChatLastTimeCleared() int64 {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	return u.spammedChatLastTimeCleared
}

// SetSpammedChatLastTimeCleared Sets the time the user's chat spam rate was cleared
func (u *User) SetSpammedChatLastTimeCleared(time int64) {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	u.spammedChatLastTimeCleared = time
}

// GetMultiplayerGameId Gets the id of the multiplayer game the user is currently inside of (if any)
func (u *User) GetMultiplayerGameId() int {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	return u.multiplayerGameId
}

// SetMultiplayerGameId Sets the id of the multiplayer game if the user is inside of one
func (u *User) SetMultiplayerGameId(id int) {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	u.multiplayerGameId = id
}

// IsMuted Returns if the user is muted
func (u *User) IsMuted() bool {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	return u.Info.MuteEndTime > time.Now().UnixMilli()
}

// MuteUser Mutes a user for a specified duration
func (u *User) MuteUser(duration time.Duration) error {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	endTime := time.Now().UnixMilli() + duration.Milliseconds()

	err := db.MuteUser(u.Info.Id, endTime)

	if err != nil {
		log.Printf("Failed to update user mute time: %v\n", err)
		return err
	}

	u.Info.MuteEndTime = endTime
	return nil
}

// GetSpectators Returns the people who are currently spectating this user
func (u *User) GetSpectators() []*User {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	return u.spectators
}

// GetSpectating Returns the people the user is currently spectating
func (u *User) GetSpectating() []*User {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	return u.spectating
}

// AddSpectator Adds a person to the list of spectators
func (u *User) AddSpectator(spectator *User) {
	clientStatus := u.GetClientStatus()

	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	if u.Info.Id == spectator.Info.Id {
		return
	}

	if utils.Includes(u.spectators, spectator) {
		return
	}

	u.spectators = append(u.spectators, spectator)
	SendPacketToUser(packets.NewServerUserInfo([]*objects.PacketUser{spectator.SerializeForPacket()}), u)
	SendPacketToUser(packets.NewServerSpectatorJoined(spectator.Info.Id), u)

	spectator.spectating = append(spectator.spectating, u)
	SendPacketToUser(packets.NewServerUserStatusSingle(u.Info.Id, clientStatus), spectator)
	SendPacketToUser(packets.NewServerStartSpectatePlayer(u.Info.Id), spectator)
	if u.status.Status == objects.ClientStatusPLaying {
		SendPacketToUser(packets.NewServerSpectatorReplayFrames(u.Info.Id, packets.SpectatorFrameNewSong, 0, nil), spectator)
	}

	// In the event that the user is already being spectated, dump all the previous frames to them so they can join in the middle.
	for _, frame := range u.frames {
		SendPacketToUser(packets.NewServerSpectatorReplayFrames(u.Info.Id, frame.Status, frame.AudioTime, frame.Frames), spectator)
	}

	runSpectatorHandlers(spectatorAddedHandlers, u, spectator)
}

// RemoveSpectator Removes a person from their list of spectators
func (u *User) RemoveSpectator(spectator *User) {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	u.spectators = utils.Filter(u.spectators, func(x *User) bool { return x != spectator })
	SendPacketToUser(packets.NewServerSpectatorLeft(spectator.Info.Id), u)

	spectator.spectating = utils.Filter(spectator.spectating, func(x *User) bool { return x != u })
	SendPacketToUser(packets.NewServerStopSpectatePlayer(u.Info.Id), spectator)

	runSpectatorHandlers(spectatorLeftHandlers, u, spectator)
}

// StopSpectatingAll Stops spectating every user that they are currently spectating
func (u *User) StopSpectatingAll() {
	for _, user := range u.GetSpectating() {
		user.RemoveSpectator(u)
	}
}

func (u *User) ClearReplayFrames() {
	u.frames = []*packets.ClientSpectatorReplayFrames{}
}

// HandleNewSpectatorFrames Handles incoming replay frames
func (u *User) HandleNewSpectatorFrames(packet *packets.ClientSpectatorReplayFrames) {
	u.Mutex.Lock()

	if u.frames == nil {
		u.ClearReplayFrames()
	}

	switch packet.Status {
	case packets.SpectatorFrameNewSong, packets.SpectatorFrameSelectingSong:
		u.ClearReplayFrames()
	default:
		u.frames = append(u.frames, packet)
	}

	u.Mutex.Unlock()

	for _, spectator := range u.GetSpectators() {
		if packet.Status == packets.SpectatorFrameNewSong || packet.Status == packets.SpectatorFrameSelectingSong {
			SendPacketToUser(packets.NewServerUserStatusSingle(u.Info.Id, u.status), spectator)
		}

		SendPacketToUser(packets.NewServerSpectatorReplayFrames(u.Info.Id, packet.Status, packet.AudioTime, packet.Frames), spectator)
	}
}

// SendClientStatusToSpectators Sends an updated user client status to all spectators
func (u *User) SendClientStatusToSpectators() {
	for _, spectator := range u.GetSpectators() {
		SendPacketToUser(packets.NewServerUserStatusSingle(u.Info.Id, u.status), spectator)
	}
}

// SerializeForPacket Serializes the user to be used in a packet
func (u *User) SerializeForPacket() *objects.PacketUser {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

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
