package packets

import "example.com/Quaver/Z/common"

type ServerGameAllowedModesChanged struct {
	Packet
	Modes []common.Mode `json:"m"`
}

func NewServerGameAllowedModesChanged(modes []common.Mode) *ServerGameAllowedModesChanged {
	return &ServerGameAllowedModesChanged{
		Packet: Packet{Id: PacketIdServerGameAllowedModesChanged},
		Modes:  modes,
	}
}
