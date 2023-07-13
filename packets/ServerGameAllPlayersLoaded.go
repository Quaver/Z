package packets

type ServerGameAllPlayersLoaded struct {
	Packet
}

func NewServerGameAllPlayersLoaded() *ServerGameAllPlayersLoaded {
	return &ServerGameAllPlayersLoaded{
		Packet{Id: PacketIdServerAllPlayersLoaded},
	}
}
