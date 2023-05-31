package packets

import (
	"example.com/Quaver/Z/objects"
)

type ServerUserConnected struct {
	Packet
	User *objects.PacketUser `json:"u"`
}

func NewServerUserConnected(user *objects.PacketUser) *ServerUserConnected {
	return &ServerUserConnected{
		Packet: Packet{Id: PacketIdServerUserConnected},
		User:   user,
	}
}
