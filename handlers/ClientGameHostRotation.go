package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client wants to enable/disable host rotation for their multiplayer game
func handleClientGameHostRotation(user *sessions.User, packet *packets.ClientGameHostRotation) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.RunLocked(func() {
		game.SetHostRotation(user, packet.Enabled)
	})
}
