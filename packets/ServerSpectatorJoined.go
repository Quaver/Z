package packets

type ServerSpectatorJoined struct {
	Packet
	UserId int `json:"u"`
}

func NewServerSpectatorJoined(userId int) *ServerSpectatorJoined {
	return &ServerSpectatorJoined{
		Packet: Packet{Id: PacketIdServerSpectatorJoined},
		UserId: userId,
	}
}
