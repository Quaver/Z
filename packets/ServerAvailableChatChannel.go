package packets

type ServerAvailableChatChannel struct {
	Packet
	Name        string `json:"n"`
	Description string `json:"d"`
	LimitedChat bool   `json:"lc"`
}

func NewServerAvailableChatChannel(name string, description string, limitedChat bool) *ServerAvailableChatChannel {
	return &ServerAvailableChatChannel{
		Packet:      Packet{Id: PacketIdServerAvailableChatChannel},
		Name:        name,
		Description: description,
		LimitedChat: limitedChat,
	}
}
