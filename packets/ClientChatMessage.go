package packets

type ClientChatMessage struct {
	Packet
	Receiver string `json:"to"`
	Message  string `json:"m"`
}
