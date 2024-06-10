package sessions

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"github.com/gobwas/ws/wsutil"
	"net"
)

// SendPacketToConnection Sends a packet to a given connection
func SendPacketToConnection(data interface{}, conn net.Conn) error {
	if conn == nil {
		return errors.New("no connection")
	}

	user := GetUserByConnection(conn)
	if user != nil {
		user.ConnMutex.Lock()
		defer user.ConnMutex.Unlock()
	}

	j, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = wsutil.WriteServerText(conn, j)

	return err
}

// SendPacketToUser Sends a packet to a given user
func SendPacketToUser(data interface{}, user *User) {
	user.PacketChannel <- data
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
