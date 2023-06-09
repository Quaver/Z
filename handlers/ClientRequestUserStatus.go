package handlers

import (
	"example.com/Quaver/Z/objects"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when the client is requesting user statuses
func handleClientRequestUserStatus(user *sessions.User, packet *packets.ClientRequestUserStatus) {
	if packet == nil {
		return
	}

	var statuses packets.ClientStatus = map[int]*objects.ClientStatus{}

	for _, userId := range packet.UserIds {
		user := sessions.GetUserById(userId)

		if user == nil {
			continue
		}

		statuses[userId] = user.GetClientStatus()
	}

	sessions.SendPacketToUser(packets.NewServerUserStatus(statuses), user)
}
