package packets

import "example.com/Quaver/Z/sessions"

type ServerUserInfo struct {
	Packet
	Users []*sessions.PacketUser `json:"u"`
}

func NewServerUserInfo(users []*sessions.PacketUser) *ServerUserInfo {
	return &ServerUserInfo{
		Packet: Packet{Id: PacketIdServerUserInfo},
		Users:  users,
	}
}
