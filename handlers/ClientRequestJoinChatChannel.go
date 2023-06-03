package handlers

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the user wants to join a chat channel
func handleClientRequestJoinChatChannel(user *sessions.User, packet *packets.ClientRequestJoinChatChannel) {
	if packet == nil {
		return
	}

	channel := chat.GetChannelByName(packet.Channel)

	if channel == nil {
		return
	}

	channel.AddUser(user)
}
