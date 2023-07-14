package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client states that their screen has loaded in-game
func handleClientGameScreenLoaded(user *sessions.User, packet *packets.ClientGameScreenLoaded) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.RunLocked(func() {
		game.SetPlayerScreenLoaded(user.Info.Id)
	})
}
