package packets

type ServerGameMapsetShared struct {
	Packet
	IsShared bool `json:"s"`
}

func NewServerGameMapsetShared(isShared bool) *ServerGameMapsetShared {
	return &ServerGameMapsetShared{
		Packet:   Packet{Id: PacketIdServerGameMapsetShared},
		IsShared: isShared,
	}
}
