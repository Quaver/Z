package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client accepts a game invite
func handleClientGameAcceptInvite(user *sessions.User, packet *packets.ClientGameAcceptInvite) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameByIdString(packet.MatchId)

	if game == nil {
		sessions.SendPacketToUser(packets.NewServerJoinGameFailed(packets.JoinGameErrorMatchNoExists), user)
		return
	}

	game.RunLocked(func() {
		game.AddPlayer(user.Info.Id, "")

		// User joined game successfully
		if multiplayer.GetGameById(user.GetMultiplayerGameId()) != nil {
			addUserToGameChat(user, game)
		}
	})
}
