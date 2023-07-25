package spectator

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
)

type User struct {
	*sessions.User
	spectators []*User // People who are currently spectating this user
	spectating []*User // People who the user is currently spectating
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

	if utils.Includes(u.spectators, spectator) {
		return
	}

	u.spectators = append(u.spectators, spectator)
	spectator.spectating = append(spectator.spectating, u)

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
	spectator.spectating = utils.Filter(spectator.spectating, func(x *User) bool { return x != u })

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
