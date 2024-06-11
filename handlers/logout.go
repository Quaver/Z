package handlers

import (
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
	"log"
	"net"
	"time"
)

func HandleLogout(conn net.Conn, temporary bool) error {
	user := sessions.GetUserByConnection(conn)

	if user != nil {
		if temporary {
			user.SetLastTemporaryDisconnectionTimestamp(time.Now().UnixMilli())
			log.Printf("[%v #%v] Logged out temporarily.\n", user.Info.Username, user.Info.Id)
		} else {
			game := multiplayer.GetGameById(user.GetMultiplayerGameId())

			if game != nil {
				game.RemovePlayer(user.Info.Id)
			}

			chat.RemoveUserFromAllChannels(user)
			multiplayer.RemoveUserFromLobby(user)

			sessions.SendPacketToAllUsers(packets.NewServerUserDisconnected(user.Info.Id))

			err := sessions.RemoveUser(user)

			if err != nil {
				log.Printf("[%v %v] Error while logging out user - %v\n", user.Info.Username, user.Info.Id, err)
			}

			log.Printf("[%v #%v] Logged out (%v users online).\n", user.Info.Username, user.Info.Id, sessions.GetOnlineUserCount())
		}
	}

	utils.CloseConnection(conn)
	return nil
}
