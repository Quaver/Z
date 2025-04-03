package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client wants to enable/disable preview for their multiplayer game
func handleClientGameEnablePreview(user *sessions.User, packet *packets.ClientGameEnablePreview) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.RunLocked(func() {
		game.SetEnablePreview(user, packet.Enabled)
	})
}
