package handlers

import (
	"example.com/Quaver/Z/chat"
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

	game.RunLocked(func() {
		game.AddPlayer(user.Info.Id, packet.Password)
		addUserToGameChat(user, game)
	})
}

// Adds a user to a game chat
func addUserToGameChat(user *sessions.User, game *multiplayer.Game) {
	chatChannel := chat.GetMultiplayerChannel(game.Data.GameId)

	if chatChannel == nil {
		return
	}

	chatChannel.AddUser(user)
}

// Removes a user from a game chat
func removeUserFromGameChat(user *sessions.User, game *multiplayer.Game) {
	chatChannel := chat.GetMultiplayerChannel(game.Data.GameId)

	if chatChannel == nil {
		return
	}

	chatChannel.RemoveUser(user)
}
