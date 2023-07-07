package packets

import "time"

type ServerGameStartCountdown struct {
	Packet
	Timestamp int64 `json:"t"`
}

func NewServerGameStartCountdown() *ServerGameStartCountdown {
	return &ServerGameStartCountdown{
		Packet:    Packet{PacketIdServerGameStartCountdown},
		Timestamp: time.Now().UnixMilli() + 5000,
	}
}
