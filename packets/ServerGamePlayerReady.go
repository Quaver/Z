package packets

type ServerGamePlayerReady struct {
	Packet
	UserId int `json:"u"`
}

func NewServerGamePlayerReady(userId int) *ServerGamePlayerReady {
	return &ServerGamePlayerReady{
		Packet: Packet{Id: PacketIdServerGamePlayerReady},
		UserId: userId,
	}
}
