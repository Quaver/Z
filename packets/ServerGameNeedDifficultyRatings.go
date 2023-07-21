package packets

type ServerGameNeedDifficultyRatings struct {
	Packet
	Md5            string `json:"md5"`
	AlternativeMd5 string `json:"amd5"`
	Needs          bool   `json:"n"`
}

func NewServerGameNeedDifficultyRatings(md5 string, alternativeMd5 string, needs bool) *ServerGameNeedDifficultyRatings {
	return &ServerGameNeedDifficultyRatings{
		Packet:         Packet{Id: PacketIdServerGameNeedDifficultyRatings},
		Md5:            md5,
		AlternativeMd5: alternativeMd5,
		Needs:          needs,
	}
}
