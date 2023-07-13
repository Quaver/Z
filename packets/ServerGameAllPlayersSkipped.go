package packets

type ServerGameAllPlayersSkipped struct {
	Packet
}

func NewServerGameAllPlayersSkipped() *ServerGameAllPlayersSkipped {
	return &ServerGameAllPlayersSkipped{
		Packet{Id: PacketIdServerGameAllPlayersSkipped},
	}
}
