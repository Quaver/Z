package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client states that they don't have the current map in multiplayer
func handleClientGamePlayerNoMap(user *sessions.User, packet *packets.ClientGamePlayerNoMap) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.RunLocked(func() {
		game.SetPlayerDoesntHaveMap(user.Info.Id)
	})
}
