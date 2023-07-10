package packets

type ServerGameKicked struct {
	Packet
}

func NewServerGameKicked() *ServerGameKicked {
	return &ServerGameKicked{Packet{Id: PacketIdServerGameKicked}}
}
