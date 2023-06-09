package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client requests to leave the multiplayer lobby
func handleClientLobbyLeave(user *sessions.User, packet *packets.ClientLobbyLeave) {
	if packet == nil {
		return
	}

	multiplayer.RemoveUserFromLobby(user)
}
