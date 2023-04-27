package packets

type ServerNotification struct {
	Packet
	Type    ServerNotificationType `json:"t"`
	Content string                 `json:"c"`
}

type ServerNotificationType int

const (
	ServerNotificationTypeError = iota
	ServerNotificationTypeSuccess
	ServerNotificationTypeInfo
)

func NewServerNotification(notificationType ServerNotificationType, content string) *ServerNotification {
	return &ServerNotification{
		Packet:  Packet{Id: PacketIdServerNotification},
		Type:    notificationType,
		Content: content,
	}
}

func NewServerNotificationError(content string) *ServerNotification {
	return NewServerNotification(ServerNotificationTypeError, content)
}

func NewServerNotificationSuccess(content string) *ServerNotification {
	return NewServerNotification(ServerNotificationTypeSuccess, content)
}

func NewServerNotificationInfo(content string) *ServerNotification {
	return NewServerNotification(ServerNotificationTypeInfo, content)
}
