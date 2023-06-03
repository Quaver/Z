package packets

type ClientRequestJoinChatChannel struct {
	Packet
	Channel string `json:"c"`
}
