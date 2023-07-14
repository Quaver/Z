package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client is requesting to change the map for their MP game.
func handleClientChangeGameMap(user *sessions.User, packet *packets.ClientChangeGameMap) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.RunLocked(func() {
		game.ChangeMap(user, packet)
	})
}
