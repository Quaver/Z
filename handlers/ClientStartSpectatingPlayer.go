package handlers

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client requests to start spectating a player.
func handleClientStartSpectatingPlayer(user *sessions.User, packet *packets.ClientStartSpectatingPlayer) {
	if packet == nil {
		return
	}

	user.StopSpectatingAll()

	spectatee := sessions.GetUserById(packet.UserId)

	if spectatee == nil {
		return
	}

	spectatee.AddSpectator(user)
}
