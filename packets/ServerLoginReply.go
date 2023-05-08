package packets

import (
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/sessions"
)

type ServerLoginReply struct {
	Packet
	User         *sessions.PacketUser  `json:"u"`
	SessionToken string                `json:"t"`
	Stats        []*db.PacketUserStats `json:"s"`
}

func NewServerLoginReply(user *sessions.User) *ServerLoginReply {
	statSlice := make([]*db.PacketUserStats, 0)

	for _, value := range user.GetStats() {
		statSlice = append(statSlice, value.SerializeForPacket())
	}

	return &ServerLoginReply{
		Packet:       Packet{Id: PacketIdServerLoginReply},
		User:         user.SerializeForPacket(),
		SessionToken: user.GetToken(),
		Stats:        statSlice,
	}
}
