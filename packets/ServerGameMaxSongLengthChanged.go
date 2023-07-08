package packets

type ServerGameMaxSongLengthChanged struct {
	Packet
	Seconds int `json:"s"`
}

func NewServerGameMaxSongLengthChanged(seconds int) *ServerGameMaxSongLengthChanged {
	return &ServerGameMaxSongLengthChanged{
		Packet:  Packet{Id: PacketIdServerGameMaxSongLengthChanged},
		Seconds: seconds,
	}
}
