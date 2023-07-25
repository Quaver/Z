package packets

type ServerStartSpectatePlayer struct {
	Packet
	UserId int `json:"u"`
}

func NewServerStartSpectatePlayer(userId int) *ServerStartSpectatePlayer {
	return &ServerStartSpectatePlayer{
		Packet: Packet{Id: PacketIdServerStartSpectatePlayer},
		UserId: userId,
	}
}
