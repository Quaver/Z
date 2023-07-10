package packets

type ServerGameHostRotation struct {
	Packet
	Enabled bool `json:"h"`
}

func NewServerGameHostRotation(enabled bool) *ServerGameHostRotation {
	return &ServerGameHostRotation{
		Packet:  Packet{Id: PacketIdServerGameHostRotationChanged},
		Enabled: enabled,
	}
}
