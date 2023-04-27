package sessions

import (
	"encoding/json"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
)

// SendPacketToConnection Sends a packet to a given connection
func SendPacketToConnection(data interface{}, conn net.Conn) error {
	j, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = wsutil.WriteServerText(conn, j)

	if err != nil {
		return err
	}

	log.Printf("SENT - %v", string(j))
	return nil
}

// SendPacketToUser Sends a packet to a given user
func SendPacketToUser(data interface{}, user *User) error {
	err := SendPacketToConnection(data, user.Conn)

	if err != nil {
		return err
	}

	return err
}

// SendPacketToUsers Sends a packet to a list of users
func SendPacketToUsers(data interface{}, users ...*User) error {
	for _, user := range users {
		err := SendPacketToUser(data, user)

		if err != nil {
			return err
		}
	}

	return nil
}
