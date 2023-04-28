package handlers

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
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

		sessions.SendPacketToAllUsers(packets.NewServerUserDisconnected(user.Info.Id))
		log.Printf("[%v #%v] Logged out (%v users online).\n", user.Info.Username, user.Info.Id, sessions.GetOnlineUserCount())
	}

	return nil
}
