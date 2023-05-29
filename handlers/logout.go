package handlers

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
	"log"
	"net"
)

func HandleLogout(conn net.Conn) error {
	user := sessions.GetUserByConnection(conn)

	if user != nil {
		err := sessions.RemoveUser(user)

		if err != nil {
			return err
		}

		chat.RemoveUserFromAllChannels(user)
		sessions.SendPacketToAllUsers(packets.NewServerUserDisconnected(user.Info.Id))

		log.Printf("[%v #%v] Logged out (%v users online).\n", user.Info.Username, user.Info.Id, sessions.GetOnlineUserCount())
	}

	utils.CloseConnection(conn)
	return nil
}
