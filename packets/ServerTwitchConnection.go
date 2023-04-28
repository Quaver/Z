package packets

type ServerTwitchConnection struct {
	Packet
	Connected bool    `json:"c"`
	Username  *string `json:"u,omitempty"`
}

func NewServerTwitchConnection(username string) *ServerTwitchConnection {
	return &ServerTwitchConnection{
		Packet:    Packet{Id: PacketIdServerTwitchConnection},
		Connected: username != "",
		Username:  &username,
	}
}
