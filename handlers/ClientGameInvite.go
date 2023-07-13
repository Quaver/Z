package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client wishes to invite someone to their game
func handleClientGameInvite(user *sessions.User, packet *packets.ClientGameInvite) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	invitee := sessions.GetUserById(packet.UserId)

	if invitee == nil {
		return
	}

	game.SendInvite(user, invitee)
}
