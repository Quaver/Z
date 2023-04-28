package packets

type ServerFriendsList struct {
	Packet
	UserIds []int `json:"u"`
}

func NewServerFriendsList(userIds []int) *ServerFriendsList {
	return &ServerFriendsList{
		Packet:  Packet{Id: PacketIdServerUserFriendsList},
		UserIds: userIds,
	}
}
