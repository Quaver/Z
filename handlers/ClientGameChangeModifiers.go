package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client wants to change the modifiers of a multiplayer game
func handleClientGameChangeModifiers(user *sessions.User, packet *packets.ClientGameChangeModifiers) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.SetGlobalModifiers(user, packet.Modifiers, packet.DifficultyRating)
}
