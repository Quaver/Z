package db

import "time"

type ChatMessageType int

const (
	ChatMessageTypePublic ChatMessageType = iota
	ChatMessageTypePrivate
)

// InsertPublicChatMessage Inserts a public chat message into the database
func InsertPublicChatMessage(senderId int, channel string, message string) error {
	err := InsertChatMessage(senderId, ChatMessageTypePublic, -1, channel, message)

	if err != nil {
		return err
	}

	return nil
}

// InsertPrivateChatMessage Inserts a private chat message into the database
func InsertPrivateChatMessage(senderId int, receiverId int, receiverName string, message string) error {
	err := InsertChatMessage(senderId, ChatMessageTypePrivate, receiverId, receiverName, message)

	if err != nil {
		return err
	}

	return nil
}

// InsertChatMessage Inserts a chat message into the database
func InsertChatMessage(senderId int, msgType ChatMessageType, receiverId int, channel string, message string) error {
	query := "INSERT INTO chat_messages (sender_id, type, receiver_id, channel, message, timestamp) VALUES (?, ?, ?, ?, ?, ?)"

	_, err := SQL.Exec(query, senderId, msgType, receiverId, channel, message, time.Now().UnixMilli())

	if err != nil {
		return err
	}

	return nil
}
