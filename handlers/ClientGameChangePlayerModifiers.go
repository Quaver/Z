package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client requests to change their modifiers in multiplayer
func handleClientGameChangePlayerModifiers(user *sessions.User, packet *packets.ClientGameChangePlayerModifiers) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.SetPlayerModifiers(user.Info.Id, packet.Modifiers)
}
