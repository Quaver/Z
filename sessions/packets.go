package sessions

import (
	"encoding/json"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
)

// SendPacketToConnection Sends a packet to a given connection
func SendPacketToConnection(data interface{}, conn net.Conn) {
	if conn == nil {
		return
	}

	j, err := json.Marshal(data)

	if err != nil {
		// log.Println(err)
		return
	}

	err = wsutil.WriteServerText(conn, j)

	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("SENT - %v", string(j))
}

// SendPacketToUser Sends a packet to a given user
func SendPacketToUser(data interface{}, user *User) {
	SendPacketToConnection(data, user.Conn)
	return
}

// SendPacketToUsers Sends a packet to a list of users
func SendPacketToUsers(data interface{}, users ...*User) {
	for _, user := range users {
		SendPacketToUser(data, user)
	}
}

// SendPacketToAllUsers Sends a packet to every online user
func SendPacketToAllUsers(data interface{}) {
	for _, user := range GetOnlineUsers() {
		SendPacketToUser(data, user)
	}
}
