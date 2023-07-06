package packets

type ServerGamePlayerNoMap struct {
	Packet
	UserId int `json:"uid"`
}

func NewServerGamePlayerNoMap(userId int) *ServerGamePlayerNoMap {
	return &ServerGamePlayerNoMap{
		Packet: Packet{Id: PacketIdServerGamePlayerNoMap},
		UserId: userId,
	}
}
