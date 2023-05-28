package packets

type ServerJoinedChatChannel struct {
	Packet
	Channel string `json:"c"`
}

func NewServerJoinedChatChannel(name string) *ServerJoinedChatChannel {
	return &ServerJoinedChatChannel{
		Packet:  Packet{Id: PacketIdServerJoinedChatChannel},
		Channel: name,
	}
}
