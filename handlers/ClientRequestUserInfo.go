package handlers

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when a client is requesting updated user info
func handleClientRequestUserInfo(user *sessions.User, packet *packets.ClientRequestUserInfo) {
	if packet == nil {
		return
	}

	userInfo := getSerializedUsersFromUserIds(packet.UserIds)
	sessions.SendPacketToUser(packets.NewServerUserInfo(userInfo), user)
}

func getSerializedUsersFromUserIds(userIds []int) []*sessions.PacketUser {
	var userInfo []*sessions.PacketUser

	for _, id := range userIds {
		user := sessions.GetUserById(id)

		if user == nil {
			continue
		}

		userInfo = append(userInfo, user.SerializeForPacket())
	}

	return userInfo
}
