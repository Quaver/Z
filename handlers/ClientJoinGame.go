package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when a user attempts to join a game
func handleClientJoinGame(user *sessions.User, packet *packets.ClientJoinGame) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameByIdString(packet.GameId)

	if game == nil {
		sessions.SendPacketToUser(packets.NewServerJoinGameFailed(packets.JoinGameErrorMatchNoExists), user)
		return
	}

	game.AddPlayer(user.Info.Id, packet.Password)
}
