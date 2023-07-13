package packets

type ServerGameSetReferee struct {
	Packet
	UserId int `json:"u"`
}

func NewServerGameSetReferee(userId int) *ServerGameSetReferee {
	return &ServerGameSetReferee{
		Packet: Packet{Id: PacketIdServerGameSetReferee},
		UserId: userId,
	}
}
