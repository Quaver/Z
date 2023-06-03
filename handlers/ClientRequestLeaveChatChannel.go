package handlers

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when a client requests to leave a given chat channel
func handleClientRequestLeaveChatChannel(user *sessions.User, packet *packets.ClientRequestLeaveChatChannel) {
	if packet == nil {
		return
	}

	channel := chat.GetChannelByName(packet.Channel)

	if channel == nil {
		return
	}

	channel.RemoveUser(user)
}
