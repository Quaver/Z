package packets

import "example.com/Quaver/Z/common"

type ClientGameJudgements struct {
	Packet
	Judgements   []common.Judgements `json:"j"`
	MineHitDelta int                 `json:"m"`
}
