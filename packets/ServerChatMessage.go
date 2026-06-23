package packets

import "time"

type ServerChatMessage struct {
	Packet
	SenderId              int    `json:"sid"`
	SenderName            string `json:"u"`
	SenderClanTag         string `json:"sct"`
	SenderClanAccentColor string `json:"sca"`
	To                    string `json:"to"` // The #channel or player username
	Message               string `json:"m"`
	Time                  int64  `json:"ts"`
}

func NewServerChatMessage(senderId int, senderName string, senderClanTag string, senderClanAccentColor string, to string, message string) *ServerChatMessage {
	return &ServerChatMessage{
		Packet:                Packet{Id: PacketIdServerChatMessage},
		SenderId:              senderId,
		SenderName:            senderName,
		SenderClanTag:         senderClanTag,
		SenderClanAccentColor: senderClanAccentColor,
		To:                    to,
		Message:               message,
		Time:                  time.Now().UnixMilli(),
	}
}
