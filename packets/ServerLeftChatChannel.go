package packets

type ServerLeftChatChannel struct {
	Packet
	Channel string `json:"c"`
}

func NewServerLeftChatChannel(name string) *ServerLeftChatChannel {
	return &ServerLeftChatChannel{
		Packet:  Packet{Id: PacketIdServerLeftChatChannelPacket},
		Channel: name,
	}
}
