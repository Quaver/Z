package handlers

import (
	"example.com/Quaver/Z/sessions"
	"net"
	"time"
)

func HandlePong(conn net.Conn) error {
	user := sessions.GetUserByConnection(conn)

	if user != nil {
		user.SetLastWsPongTimestamp(time.Now().UnixMilli())
	}

	return nil
}
