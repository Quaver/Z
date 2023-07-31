package packets

type ServerGameAutoHost struct {
	Packet
	Enabled bool `json:"enabled"`
}

func NewServerGameAutoHost(enabled bool) *ServerGameAutoHost {
	return &ServerGameAutoHost{
		Packet:  Packet{Id: PacketIdServerGameAutoHost},
		Enabled: enabled,
	}
}
