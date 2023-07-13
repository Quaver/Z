package packets

type ServerGameEnded struct {
	Packet
}

func NewServerGameEnded() *ServerGameEnded {
	return &ServerGameEnded{
		Packet{Id: PacketIdServerGameEnded},
	}
}
