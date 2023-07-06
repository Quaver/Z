package packets

type ServerGameDisbanded struct {
	Packet
	GameId string `json:"gid"`
}

func NewServerGameDisbanded(gameId string) *ServerGameDisbanded {
	return &ServerGameDisbanded{
		Packet: Packet{Id: PacketIdServerGameDisbanded},
		GameId: gameId,
	}
}
