package packets

type ServerGameChangePlayerModifiers struct {
	Packet
	UserId    int   `json:"u"`
	Modifiers int64 `json:"m"`
}

func NewServerGameChangePlayerModifiers(userId int, mods int64) *ServerGameChangePlayerModifiers {
	return &ServerGameChangePlayerModifiers{
		Packet:    Packet{Id: PacketIdServerGamePlayerChangeModifiers},
		UserId:    userId,
		Modifiers: mods,
	}
}
