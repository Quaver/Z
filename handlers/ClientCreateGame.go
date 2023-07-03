package handlers

import (
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"log"
)

// Handles when the user wants to create a multiplayer game.
func handleClientCreateGame(user *sessions.User, packet *packets.ClientCreateGame) {
	if packet == nil || packet.Game == nil {
		return
	}

	game, err := multiplayer.NewGame(packet.Game)
	game.Data.HostId = user.Info.Id

	if err != nil {
		log.Printf("Error creating multiplayer game: %v\n", err)
		return
	}

	multiplayer.AddGameToLobby(game)
}
