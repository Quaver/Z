package handlers

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Called when the client requests to log out.
func handleClientLogout(user *sessions.User, packet *packets.ClientLogout) {
	if packet == nil {
		return
	}

	err := HandleLogout(user.Conn, false)

	if err != nil {
		return
	}
}
