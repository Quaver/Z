package handlers

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client requests to leave a game
func handleClientLeaveGame(user *sessions.User, packet *packets.ClientLeaveGame) {
	if packet == nil {
		return
	}

	game := multiplayer.GetGameById(user.GetMultiplayerGameId())

	if game == nil {
		return
	}

	game.RunLocked(func() {
		game.RemovePlayer(user.Info.Id)
		removeUserFromGameChat(user, game)

		if len(game.Data.PlayerIds) == 0 {
			chat.RemoveMultiplayerChannel(game.Data.GameId)
		}
	})
}
