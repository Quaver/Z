package packets

import "example.com/Quaver/Z/objects"

type ServerGameChangeFreeMod struct {
	Packet
	Type objects.MultiplayerGameFreeMod `json:"fm"`
}

func NewServerGameChangeFreeMod(freeMod objects.MultiplayerGameFreeMod) *ServerGameChangeFreeMod {
	return &ServerGameChangeFreeMod{
		Packet: Packet{Id: PacketIdServerGameFreeModTypeChanged},
		Type:   freeMod,
	}
}
