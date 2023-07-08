package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client is stating that they are/aren't selecting a map
func handleClientGameHostSelectingMap(user *sessions.User, packet *packets.ClientGameHostSelectingMap) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.SetHostSelectingMap(user, packet.IsSelecting, true, true)
}
