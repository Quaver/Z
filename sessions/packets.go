package sessions

import (
	"encoding/json"
	"example.com/Quaver/Z/common"
	"github.com/gobwas/ws/wsutil"
	"net"
	"sync"
	"time"
)

var (
	packetQueue map[*User][]interface{} // Packets that are queued to be sent
	packetMutex *sync.Mutex             // Locks packetQueue
)

func init() {
	StartSendingPackets()
}

// StartSendingPackets Begins sending out packets to users
func StartSendingPackets() {
	packetQueue = map[*User][]interface{}{}
	packetMutex = &sync.Mutex{}

	go func() {
		for {
			packetMutex.Lock()

			for user, packets := range packetQueue {
				// The user's session has been fully closed, so we can stop sending packets to them.
				if user.SessionClosed {
					delete(packetQueue, user)
					continue
				}

				var packetsFailed []interface{}

				// Try to send packets, and update the packets that failed to send to try again
				for _, packet := range packets {
					err := SendPacketToConnection(packet, user.Conn)

					if err != nil {
						packetsFailed = append(packetsFailed, packet)
					}
				}

				packetQueue[user] = packetsFailed
			}

			packetMutex.Unlock()
			time.Sleep(time.Millisecond * 10)
		}
	}()
}

// SendPacketToConnection Sends a packet to a given connection
func SendPacketToConnection(data interface{}, conn net.Conn) (err error) {
	j, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = wsutil.WriteServerText(conn, j)
	return err
}

// SendPacketToUser Sends a packet to a given user
func SendPacketToUser(data interface{}, user *User) {
	packetMutex.Lock()
	defer packetMutex.Unlock()

	// Never send packets to bot users.
	if common.HasUserGroup(user.Info.UserGroups, common.UserGroupBot) {
		return
	}

	packetQueue[user] = append(packetQueue[user], data)
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
