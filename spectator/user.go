package spectator

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
)

type User struct {
	*sessions.User
	spectators []*User                                // People who are currently spectating this user
	spectating []*User                                // People who the user is currently spectating
	frames     []*packets.ClientSpectatorReplayFrames // The replay frames for the user's current play session
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
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	if u.Info.Id == spectator.Info.Id {
		return
	}

	if utils.Includes(u.spectators, spectator) {
		return
	}

	u.spectators = append(u.spectators, spectator)
	sessions.SendPacketToUser(packets.NewServerSpectatorJoined(spectator.Info.Id), u.User)

	spectator.spectating = append(spectator.spectating, u)
	sessions.SendPacketToUser(packets.NewServerUserStatusSingle(u.Info.Id, u.GetClientStatus()), spectator.User)
	sessions.SendPacketToUser(packets.NewServerStartSpectatePlayer(u.Info.Id), spectator.User)

	// In the event that the user is already being spectated, dump all the previous frames to them so they can join in the middle.
	for _, frame := range u.frames {
		sessions.SendPacketToUser(packets.NewServerSpectatorReplayFrames(u.Info.Id, frame.Status, frame.AudioTime, frame.Frames), spectator.User)
	}

	// Create spectator channel, and add the users to it
	if len(u.spectators) == 1 {
		channel := chat.AddSpectatorChannel(u.Info.Id)
		channel.AddUser(u.User)
		channel.AddUser(spectator.User)
		return
	}

	channel := chat.GetSpectatorChannel(u.Info.Id)

	if channel == nil {
		return
	}

	channel.AddUser(spectator.User)
}

// RemoveSpectator Removes a person from their list of spectators
func (u *User) RemoveSpectator(spectator *User) {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	u.spectators = utils.Filter(u.spectators, func(x *User) bool { return x != spectator })
	sessions.SendPacketToUser(packets.NewServerSpectatorLeft(spectator.Info.Id), u.User)

	spectator.spectating = utils.Filter(spectator.spectating, func(x *User) bool { return x != u })
	sessions.SendPacketToUser(packets.NewServerStopSpectatePlayer(u.Info.Id), spectator.User)

	channel := chat.GetSpectatorChannel(u.Info.Id)

	if channel != nil {
		channel.RemoveUser(spectator.User)

		// Remove spectator channel now that there are no longer any spectators.
		if len(u.spectators) == 0 {
			channel.RemoveUser(u.User)
			chat.RemoveSpectatorChannel(u.Info.Id)
		}
	}
}

// StopSpectatingAll Stops spectating every user that they are currently spectating
func (u *User) StopSpectatingAll() {
	for _, user := range u.GetSpectating() {
		user.RemoveSpectator(u)
	}
}

// HandleIncomingFrames Handles incoming replay frames
func (u *User) HandleIncomingFrames(packet *packets.ClientSpectatorReplayFrames) {
	u.Mutex.Lock()

	if u.frames == nil {
		u.frames = []*packets.ClientSpectatorReplayFrames{}
	}

	switch packet.Status {
	case packets.SpectatorFrameNewSong, packets.SpectatorFrameSelectingSong:
		u.frames = []*packets.ClientSpectatorReplayFrames{}
	default:
		u.frames = append(u.frames, packet)
	}

	u.Mutex.Unlock()

	for _, spectator := range u.GetSpectators() {
		sessions.SendPacketToUser(packets.NewServerUserStatusSingle(u.Info.Id, u.GetClientStatus()), spectator.User)
		sessions.SendPacketToUser(packets.NewServerSpectatorReplayFrames(u.Info.Id, packet.Status, packet.AudioTime, packet.Frames), spectator.User)
	}
}
