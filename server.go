package main

import (
	"example.com/Quaver/Z/handlers"
	"example.com/Quaver/Z/utils"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"net/http"
	"strings"
)

type Server struct {
	// The port the server is running on
	Port int

	// If the server is currently started
	IsStarted bool
}

// NewServer Creates and returns a new server object.
func NewServer(port int) *Server {
	s := Server{
		Port: port,
	}

	return &s
}

// Start Starts running the server
func (s *Server) Start() {
	if s.IsStarted {
		log.Println("Server is already started. Cannot start again!")
		return
	}

	log.Printf("Starting server on port: %v\n", s.Port)

	err := http.ListenAndServe(fmt.Sprintf(":%v", s.Port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP Connection
		conn, _, _, err := ws.UpgradeHTTP(r, w)

		if err != nil {
			log.Println(err)
			return
		}

		// Login users if necessary
		if strings.Contains(r.RequestURI, "/?login=") {
			err := handlers.HandleLogin(conn, r)

			if err != nil {
				log.Println(err)
				utils.CloseConnection(conn)
				return
			}
		}

		// Handle various connection events
		go func() {
			defer conn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(conn)

				if err != nil {
					if strings.Contains(err.Error(), "use of closed network connection") {
						return
					}

					continue
				}

				switch op {
				case ws.OpText:
					s.onTextMessage(conn, msg)
					break
				case ws.OpBinary:
					s.onBinaryMessage(conn, msg)
					break
				case ws.OpClose:
					s.onClose(conn)
					break
				case ws.OpPing:
					s.onPing(conn)
					break
				case ws.OpPong:
					s.onPong(conn)
					break
				}
			}
		}()
	}))

	if err != nil {
		panic(err)
	}
}

// Handles new incoming text messages
func (s *Server) onTextMessage(conn net.Conn, msg []byte) {
	log.Printf("Text Messsage: %v\n", string(msg))
}

// Handles new incoming binary messages
func (s *Server) onBinaryMessage(conn net.Conn, msg []byte) {
	log.Printf("Binary Messsage: %v\n", &msg)
}

// Handles when a connection has been closed
func (s *Server) onClose(conn net.Conn) {
	log.Printf("Connection closed: %v\n", conn.RemoteAddr())
}

// Handles when a connection pinged
func (s *Server) onPing(conn net.Conn) {
	log.Printf("Ping: %v\n", conn.RemoteAddr())
}

// Handles when a connection pinged
func (s *Server) onPong(conn net.Conn) {
	log.Printf("Pong: %v\n", conn.RemoteAddr())
}
