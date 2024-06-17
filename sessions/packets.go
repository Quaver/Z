package sessions

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
)

func SendPingToUser(user *User) error {
	user.ConnMutex.Lock()
	defer user.ConnMutex.Unlock()

	return wsutil.WriteServerMessage(user.Conn, ws.OpPing, []byte("ping"))
}

// SendPacketToConnection Sends a packet to a given connection
func SendPacketToConnection(data interface{}, conn net.Conn) {
	if conn == nil {
		return
	}

	user := GetUserByConnection(conn)
	if user != nil {
		user.ConnMutex.Lock()
		defer user.ConnMutex.Unlock()
	}

	j, err := json.Marshal(data)

	if err != nil {
		return
	}

	err = wsutil.WriteServerText(conn, j)

	if err != nil {
		return
	}
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
