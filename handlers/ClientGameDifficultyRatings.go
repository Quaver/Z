package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client sends us difficulty ratings to use
func handleClientGameDifficultyRatings(user *sessions.User, packet *packets.ClientGameDifficultyRatings) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.SetClientProvidedDifficultyRatings(packet.Md5, packet.AlternativeMd5, packet.Difficulties)
}
