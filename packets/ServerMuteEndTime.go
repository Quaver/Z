package packets

type ServerMuteEndTime struct {
	Packet
	UserId      int   `json:"u"`
	MuteEndTime int64 `json:"t"`
}

func NewServerMuteEndTime(userId int, muteEndTime int64) *ServerMuteEndTime {
	return &ServerMuteEndTime{
		Packet:      Packet{Id: PacketIdServerMuteEndTimePacket},
		UserId:      userId,
		MuteEndTime: muteEndTime,
	}
}
