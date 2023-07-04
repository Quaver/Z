package packets

type ServerJoinGame struct {
	Packet
	GameId string `json:"gid"`
}

func NewServerJoinGame(gameId string) *ServerJoinGame {
	return &ServerJoinGame{
		Packet: Packet{Id: PacketIdServerJoinGame},
		GameId: gameId,
	}
}
