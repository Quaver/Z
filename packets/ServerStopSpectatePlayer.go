package packets

type ServerStopSpectatePlayer struct {
	Packet
	UserId int `json:"u"`
}

func NewServerStopSpectatePlayer(userId int) *ServerStopSpectatePlayer {
	return &ServerStopSpectatePlayer{
		Packet: Packet{Id: PacketIdServerStopSpectatePlayer},
		UserId: userId,
	}
}
