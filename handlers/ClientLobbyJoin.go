package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when a client requests to join the multiplayer lobby
func handleClientLobbyJoin(user *sessions.User, packet *packets.ClientLobbyJoin) {
	if packet == nil {
		return
	}

	multiplayer.AddUserToLobby(user)
}
