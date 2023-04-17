package utils

import (
	"log"
	"net"
)

// CloseConnection Closes a connection
func CloseConnection(conn net.Conn) {
	err := conn.Close()

	if err != nil {
		log.Printf("[%v]: Failed to close connection", conn.RemoteAddr())
	}
}
