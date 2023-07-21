package packets

type ServerGameNeedDifficultyRatings struct {
	Packet
	Needs bool `json:"n"`
}

func NewServerGameNeedDifficultyRatings(needs bool) *ServerGameNeedDifficultyRatings {
	return &ServerGameNeedDifficultyRatings{
		Packet: Packet{Id: PacketIdServerGameNeedDifficultyRatings},
		Needs:  needs,
	}
}
