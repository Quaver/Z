package packets

type ServerGameEnded struct {
	Packet
	Force bool `json:"force"`
}

func NewServerGameEnded(force bool) *ServerGameEnded {
	return &ServerGameEnded{
		Packet{Id: PacketIdServerGameEnded},
		force,
	}
}
