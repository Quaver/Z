package packets

type ServerGameStart struct {
	Packet
}

func NewServerGameStart() *ServerGameStart {
	return &ServerGameStart{
		Packet{Id: PacketIdServerGameStart},
	}
}
