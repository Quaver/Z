package sessions

import (
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
)

// SendPacketToConnection Sends a packet to a given connection
func SendPacketToConnection(data interface{}, conn net.Conn) {
	if conn == nil {
		return
	}

	j, err := json.Marshal(data)

	if err != nil {
		return
	}

	writer := wsutil.GetWriter(conn, ws.StateServerSide, ws.OpBinary, wsutil.DefaultWriteBuffer)

	_, err = writer.Write(j)

	if err == nil {
		err = writer.Flush()
	}

	wsutil.PutWriter(writer)

	if err != nil {
		fmt.Printf("error: %v\n", err)
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
