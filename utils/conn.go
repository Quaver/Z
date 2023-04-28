package utils

import (
	"net"
	"time"
)

// CloseConnection Closes a connection
func CloseConnection(conn net.Conn) {
	_ = conn.Close()
}

// CloseConnectionDelayed Closes the connection after a specified amount of time
func CloseConnectionDelayed(conn net.Conn) {
	time.AfterFunc(250*time.Millisecond, func() {
		CloseConnection(conn)
	})
}
