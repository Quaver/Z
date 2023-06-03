package packets

type ClientRequestLeaveChatChannel struct {
	Packet
	Channel string `json:"c"`
}
