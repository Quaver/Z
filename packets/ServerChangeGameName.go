package packets

type ServerGameNameChanged struct {
	Packet
	Name string `json:"n"`
}

func NewServerGameNameChanged(name string) *ServerGameNameChanged {
	return &ServerGameNameChanged{
		Packet: Packet{Id: PacketIdServerGameNameChanged},
		Name:   name,
	}
}
