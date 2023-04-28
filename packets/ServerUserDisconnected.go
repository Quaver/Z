package packets

type ServerUserDisconnected struct {
	Packet
	UserId int `json:"u"`
}

func NewServerUserDisconnected(userId int) *ServerUserDisconnected {
	return &ServerUserDisconnected{
		Packet: Packet{Id: PacketIdServerUserDisconnected},
		UserId: userId,
	}
}
