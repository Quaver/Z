package handlers

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when a client requests to retrieve a collection of users' stats
func handleClientRequestUserStats(user *sessions.User, packet *packets.ClientRequestUserStats) {
	if packet == nil || packet.Users == nil || len(packet.Users) == 0 {
		return
	}

	var statsObj = map[int]map[common.Mode]*db.PacketUserStats{}

	for _, packetUser := range packet.Users {
		u := sessions.GetUserById(packetUser)

		if u == nil {
			continue
		}

		stats := u.GetStats()

		statsObj[packetUser] = map[common.Mode]*db.PacketUserStats{}

		for i, stat := range stats {
			statsObj[packetUser][i] = stat.SerializeForPacket()
		}
	}

	sessions.SendPacketToUser(packets.NewServerUserStats(statsObj), user)
}
