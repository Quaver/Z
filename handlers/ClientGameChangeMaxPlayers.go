package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client requests to change the max player count of a multiplayer game
func handleClientGameChangeMaxPlayers(user *sessions.User, packet *packets.ClientGameChangeMaxPlayers) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.SetMaxPlayerCount(user, packet.Count)
}
