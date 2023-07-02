package handlers

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"fmt"
)

// Handles when the user wants to create a multiplayer game.
func handleClientCreateGame(user *sessions.User, packet *packets.ClientCreateGame) {
	if packet == nil {
		return
	}

	fmt.Println(packet)
}
