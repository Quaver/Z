package packets

import (
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/objects"
)

type ServerLoginReply struct {
	Packet
	User         *objects.PacketUser   `json:"u"`
	SessionToken string                `json:"t"`
	Stats        []*db.PacketUserStats `json:"s"`
}

func NewServerLoginReply(user *objects.PacketUser, stats []*db.PacketUserStats, token string) *ServerLoginReply {
	return &ServerLoginReply{
		Packet:       Packet{Id: PacketIdServerLoginReply},
		User:         user,
		SessionToken: token,
		Stats:        stats,
	}
}
