package packets

type ServerAvailableChatChannel struct {
	Packet
	Name        string `json:"n"`
	Description string `json:"d"`
}

func NewServerAvailableChatChannel(name string, description string) *ServerAvailableChatChannel {
	return &ServerAvailableChatChannel{
		Packet:      Packet{Id: PacketIdServerAvailableChatChannel},
		Name:        name,
		Description: description,
	}
}
