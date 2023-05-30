package handlers

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"log"
)

// Handles when a user sends a chat message
func handleClientChatMessage(user *sessions.User, packet *packets.ClientChatMessage) {
	if packet == nil {
		return
	}

	log.Println(packet)
}
