package handlers

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client wishes to stop spectating the player
func handleClientStopSpectatingPlayer(user *sessions.User, packet *packets.ClientStopSpectatingPlayer) {
	if packet == nil {
		return
	}

	user.StopSpectatingAll()
}
