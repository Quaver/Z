package utils

import (
	"log"
	"net"
	"time"
)

// CloseConnection Closes a connection
func CloseConnection(conn net.Conn) {
	err := conn.Close()

	if err != nil {
		log.Printf("[%v]: Failed to close connection", conn.RemoteAddr())
	}
}

// CloseConnectionDelayed Closes the connection after a specified amount of time
func CloseConnectionDelayed(conn net.Conn) {
	time.AfterFunc(250*time.Millisecond, func() {
		CloseConnection(conn)
	})
}
