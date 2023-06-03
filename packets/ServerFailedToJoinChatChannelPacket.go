package packets

type ServerFailedToJoinChatChannel struct {
	Packet
	Channel string `json:"c"`
}

func NewServerFailedToJoinChatChannel(channel string) *ServerFailedToJoinChatChannel {
	return &ServerFailedToJoinChatChannel{
		Packet:  Packet{Id: PacketIdServerFailedToJoinChannelPacket},
		Channel: channel,
	}
}
