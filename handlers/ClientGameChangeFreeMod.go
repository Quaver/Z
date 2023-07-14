package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client wants to change the free mod type in a multiplayer game
func handleClientGameChangeFreeMod(user *sessions.User, packet *packets.ClientGameFreeModTypeChanged) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.RunLocked(func() {
		game.SetFreeMod(user, packet.Type)
	})
}
