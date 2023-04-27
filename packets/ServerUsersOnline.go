package packets

type ServerUsersOnline struct {
	Packet
	UserIds []int `json:"u"`
}

func NewServerUsersOnline(userIds []int) *ServerUsersOnline {
	return &ServerUsersOnline{
		Packet:  Packet{Id: PacketIdServerUsersOnline},
		UserIds: userIds,
	}
}
