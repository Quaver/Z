package sessions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gobwas/ws/wsutil"
	"net"
)

// SendPacketToConnection Sends a packet to a given connection
func SendPacketToConnection(data interface{}, conn net.Conn) (err error) {
	if conn == nil {
		return errors.New("no connection")
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = errors.New("unknown panic")
			}
		}
	}()

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
	return
}

// SendPacketToUser Sends a packet to a given user
func SendPacketToUser(data interface{}, user *User) {
	user.PacketChannel.In <- data
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
