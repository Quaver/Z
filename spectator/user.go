package spectator

import (
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

	// Create spectator chat channel
	if len(u.spectators) == 1 {

	}

	// Add user to spectator channel
}

// RemoveSpectator Removes a person from their list of spectators
func (u *User) RemoveSpectator(spectator *User) {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	u.spectators = utils.Filter(u.spectators, func(x *User) bool { return x != spectator })
	spectator.spectating = utils.Filter(spectator.spectating, func(x *User) bool { return x != u })

	// Remove user from spectator chat channel
	// Remove spectator chat channel
	if len(u.spectators) == 0 {

	}
}

// StopSpectatingAll Stops spectating every user that they are currently spectating
func (u *User) StopSpectatingAll() {
	for _, user := range u.GetSpectating() {
		user.RemoveSpectator(u)
	}
}
