package packets

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
)

type ServerUserStats struct {
	Packet
	Stats map[int]map[common.Mode]*db.PacketUserStats `json:"u"`
}

func NewServerUserStats(stats map[int]map[common.Mode]*db.PacketUserStats) *ServerUserStats {
	return &ServerUserStats{
		Packet: Packet{Id: PacketIdServerUserStats},
		Stats:  stats,
	}
}
