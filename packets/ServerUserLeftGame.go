package packets

type ServerUserLeftGame struct {
	Packet
	UserId int `json:"uid"`
}

func NewServerUserLeftGame(userId int) *ServerUserLeftGame {
	return &ServerUserLeftGame{
		Packet: Packet{Id: PacketIdServerUserLeftGame},
		UserId: userId,
	}
}
