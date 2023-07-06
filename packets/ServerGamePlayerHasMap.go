package packets

type ServerGamePlayerHasMap struct {
	Packet
	UserId int `json:"uid"`
}

func NewServerGamePlayerHasMap(userId int) *ServerGamePlayerHasMap {
	return &ServerGamePlayerHasMap{
		Packet: Packet{Id: PacketIdServerGamePlayerHasMap},
		UserId: userId,
	}
}
