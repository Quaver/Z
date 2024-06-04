package utils

import (
	"github.com/gobwas/ws"
	"net"
	"time"
)

// CloseConnection Closes a connection
func CloseConnection(conn net.Conn) {
	if conn == nil {
		return
	}

	var body = ws.NewCloseFrameBody(1000, "")
	var frame = ws.NewCloseFrame(body)
	if err := ws.WriteHeader(conn, frame.Header); err != nil {
		return
	}
	if _, err := conn.Write(body); err != nil {
		return
	}
	_ = conn.Close()
}

// CloseConnectionDelayed Closes the connection after a specified amount of time
func CloseConnectionDelayed(conn net.Conn) {
	time.AfterFunc(250*time.Millisecond, func() {
		CloseConnection(conn)
	})
}
