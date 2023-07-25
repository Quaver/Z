package packets

type ServerSpectatorLeft struct {
	Packet
	UserId int `json:"u"`
}

func NewServerSpectatorLeft(userId int) *ServerSpectatorLeft {
	return &ServerSpectatorLeft{
		Packet: Packet{Id: PacketIdServerSpectatorLeft},
		UserId: userId,
	}
}
