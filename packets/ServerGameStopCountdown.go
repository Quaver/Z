package packets

type ServerGameStopCountdown struct {
	Packet
}

func NewServerGameStopCountdown() *ServerGameStartCountdown {
	return &ServerGameStartCountdown{
		Packet:    Packet{Id: PacketIdServerGameStopCountdown},
		Timestamp: 0,
	}
}
