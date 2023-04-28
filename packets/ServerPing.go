package packets

type ServerPing struct {
	Packet
}

func NewServerPing() *ServerPing {
	return &ServerPing{
		Packet: Packet{Id: PacketIdServerPing},
	}
}
