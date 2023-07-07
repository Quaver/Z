package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when a player states that they are no longer ready in the multiplayer game
func handleClientGamePlayerNotReady(user *sessions.User, packet *packets.ClientGamePlayerNotReady) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.SetPlayerNotReady(user.Info.Id)
}
