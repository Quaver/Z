package packets

type ServerGameChangeMaxPlayers struct {
	Packet
	Count int `json:"p"`
}

func NewServerGameChangeMaxPlayers(count int) *ServerGameChangeMaxPlayers {
	return &ServerGameChangeMaxPlayers{
		Packet: Packet{Id: PacketIdServerGameMaxPlayersChanged},
		Count:  count,
	}
}
