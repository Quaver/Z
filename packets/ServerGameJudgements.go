package packets

import "example.com/Quaver/Z/common"

type ServerGameJudgements struct {
	Packet
	UserId       int                 `json:"u"`
	Judgements   []common.Judgements `json:"j"`
	MineHitDelta int                 `json:"m"`
}

func NewServerGameJudgements(userId int, judgements []common.Judgements, mineHitDelta int) *ServerGameJudgements {
	return &ServerGameJudgements{
		Packet:       Packet{Id: PacketIdServerGameJudgements},
		UserId:       userId,
		Judgements:   judgements,
		MineHitDelta: mineHitDelta,
	}
}
