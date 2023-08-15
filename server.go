package main

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/handlers"
	"example.com/Quaver/Z/multiplayer"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	// The port the server is running on
	Port int

	// If the server is currently started
	IsStarted bool
}

// NewServer Creates and returns a new server object.
func NewServer(port int) *Server {
	if port <= 0 {
		panic(fmt.Sprintf("invalid port: `%v` provided", port))
	}

	s := Server{
		Port: port,
	}

	return &s
}

// Start Starts running the server
func (s *Server) Start() {
	if s.IsStarted {
		log.Fatalln("Server is already started. Cannot start again!")
	}

	clearPreviousSessions()
	startBackgroundWorker()

	log.Printf("Starting server on port: %v\n", s.Port)

	err := http.ListenAndServe(fmt.Sprintf(":%v", s.Port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)

		if err != nil {
			log.Println(err)
			return
		}

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
					_ = s.onClose(conn)
					return
				}

				switch op {
				case ws.OpText:
					s.onTextMessage(conn, msg)
					break
				case ws.OpClose:
					err := s.onClose(conn)

					if err != nil {
						log.Println(err)
					}
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
	handlers.HandleIncomingPackets(conn, string(msg))
}

// Handles when a connection has been closed
func (s *Server) onClose(conn net.Conn) error {
	err := handlers.HandleLogout(conn)

	if err != nil {
		return err
	}

	return nil
}

// Cleans up the previous sessions (when restarting the server)
func clearPreviousSessions() {
	err := sessions.UpdateRedisOnlineUserCount()

	if err != nil {
		panic(err)
	}

	err = sessions.ClearRedisUserTokens()

	if err != nil {
		panic(err)
	}

	err = sessions.ClearRedisUserClientStatuses()

	if err != nil {
		panic(err)
	}

	err = multiplayer.ClearRedisGames()

	if err != nil {
		panic(err)
	}

	log.Println("Cleared previous redis sessions")
}

// Handles all operations that happen in the background at intervals to keep the server clean.
func startBackgroundWorker() {
	go func() {
		for {
			users := sessions.GetOnlineUsers()

			for _, user := range users {
				// Disregard bot users
				if common.HasUserGroup(user.Info.UserGroups, common.UserGroupBot) {
					continue
				}

				// Clear user's chat spam rate
				if time.Now().UnixMilli()-user.GetSpammedChatLastTimeCleared() >= 10_000 {
					user.ResetSpammedMessagesCount()
					user.SetSpammedChatLastTimeCleared(time.Now().UnixMilli())
				}

				// Ping the user periodically
				if time.Now().UnixMilli()-user.GetLastPingTimestamp() >= 40_000 {
					sessions.SendPacketToUser(packets.NewServerPing(), user)
					user.SetLastPingTimestamp()
				}

				// User hasn't responded to pings in a while, so disconnect them
				if time.Now().UnixMilli()-user.GetLastPongTimestamp() >= 120_000 {
					err := handlers.HandleLogout(user.Conn)

					if err != nil {
						log.Println(err)
					}

					log.Printf("[%v - %v] Disconnected due to being unresponsive to pings (timeout)\n", user.Info.Username, user.Info.Id)
					continue
				}
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()
}
