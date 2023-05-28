package packets

import "time"

type ServerChatMessage struct {
	Packet
	SenderId   int    `json:"sid"`
	SenderName string `json:"u"`
	To         string `json:"to"` // The #channel or player username
	Message    string `json:"m"`
	Time       int64  `json:"ts"`
}

func NewServerChatMessage(senderId int, senderName string, to string, message string) *ServerChatMessage {
	return &ServerChatMessage{
		Packet:     Packet{Id: PacketIdServerChatMessage},
		SenderId:   senderId,
		SenderName: senderName,
		To:         to,
		Message:    message,
		Time:       time.Now().UnixMilli(),
	}
}
