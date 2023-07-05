package packets

type ServerUserJoinedGame struct {
	Packet
	UserId int `json:"uid"`
}

func NewServerUserJoinedGame(userId int) *ServerUserJoinedGame {
	return &ServerUserJoinedGame{
		Packet: Packet{PacketIdServerUserJoinedGame},
		UserId: userId,
	}
}
