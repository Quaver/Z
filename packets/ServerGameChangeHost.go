package packets

type ServerGameChangeHost struct {
	Packet
	UserId int `json:"u"`
}

func NewServerGameChangeHost(id int) *ServerGameChangeHost {
	return &ServerGameChangeHost{
		Packet: Packet{Id: PacketIdServerChangeGameHost},
		UserId: id,
	}
}
