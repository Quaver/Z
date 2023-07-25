package handlers

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/spectator"
)

// Handles when the client requests to start spectating a player.
func handleClientStartSpectatingPlayer(user *sessions.User, packet *packets.ClientStartSpectatingPlayer) {
	if packet == nil {
		return
	}

	spectatee := sessions.GetUserById(packet.UserId)

	if spectatee == nil {
		return
	}

	spectUser := spectator.GetUser(user)
	spectUser.StopSpectatingAll()

	spectator.GetUser(spectatee).AddSpectator(spectUser)
}
