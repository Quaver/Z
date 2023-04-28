package packets

import "example.com/Quaver/Z/sessions"

type ServerUserConnected struct {
	Packet
	User *sessions.PacketUser `json:"u"`
}

func NewServerUserConnected(user *sessions.PacketUser) *ServerUserConnected {
	return &ServerUserConnected{
		Packet: Packet{Id: PacketIdServerUserConnected},
		User:   user,
	}
}
