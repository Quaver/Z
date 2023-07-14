package packets

import "example.com/Quaver/Z/common"

type ServerGameChangePlayerModifiers struct {
	Packet
	UserId    int         `json:"u"`
	Modifiers common.Mods `json:"m"`
}

func NewServerGameChangePlayerModifiers(userId int, mods common.Mods) *ServerGameChangePlayerModifiers {
	return &ServerGameChangePlayerModifiers{
		Packet:    Packet{Id: PacketIdServerGamePlayerChangeModifiers},
		UserId:    userId,
		Modifiers: mods,
	}
}
