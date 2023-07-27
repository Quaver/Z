package packets

type ServerSpectateMultiplayerGame struct {
	Packet
	GameId string `json:"gid"`
}

func NewServerSpectateMultiplayerGame(gameId string) *ServerSpectateMultiplayerGame {
	return &ServerSpectateMultiplayerGame{
		Packet: Packet{Id: PacketIdServerSpectateMultiplayerGame},
		GameId: gameId,
	}
}
