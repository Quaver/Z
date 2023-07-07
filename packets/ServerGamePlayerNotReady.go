package packets

type ServerGamePlayerNotReady struct {
	Packet
	UserId int `json:"u"`
}

func NewServerGamePlayerNotReady(userId int) *ServerGamePlayerReady {
	return &ServerGamePlayerReady{
		Packet: Packet{PacketIdServerGamePlayerNotReady},
		UserId: userId,
	}
}
