package packets

import (
	"example.com/Quaver/Z/objects"
)

type ServerUserInfo struct {
	Packet
	Users []*objects.PacketUser `json:"u"`
}

func NewServerUserInfo(users []*objects.PacketUser) *ServerUserInfo {
	return &ServerUserInfo{
		Packet: Packet{Id: PacketIdServerUserInfo},
		Users:  users,
	}
}
