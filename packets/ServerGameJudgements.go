package packets

import "example.com/Quaver/Z/common"

type ServerGameJudgements struct {
	Packet
	UserId     int                 `json:"u"`
	Judgements []common.Judgements `json:"j"`
}

func NewServerGameJudgements(userId int, judgements []common.Judgements) *ServerGameJudgements {
	return &ServerGameJudgements{
		Packet:     Packet{Id: PacketIdServerGameJudgements},
		UserId:     userId,
		Judgements: judgements,
	}
}
