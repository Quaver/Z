package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when a client requests to spectate a multiplayer game
func handleClientSpectateMultiplayerGame(user *sessions.User, packet *packets.ClientSpectateMultiplayerGame) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameByIdString(packet.GameId)

	if game == nil {
		return
	}

	game.RunLocked(func() {
		game.AddSpectator(user, packet.Password)
	})
}
