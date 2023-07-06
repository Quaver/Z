package packets

import "example.com/Quaver/Z/common"

type ServerGameMapChanged struct {
	Packet
	MD5                  string      `json:"md5"`
	AlternativeMD5       string      `json:"amd5"`
	MapId                int         `json:"mid"`
	MapsetId             int         `json:"msid"`
	Name                 string      `json:"map"`
	GameMode             common.Mode `json:"gm"`
	DifficultyRating     float64     `json:"d"`
	DifficultyRattingAll []float64   `json:"adr"`
	JudgementCount       int         `json:"jc"`
}

func NewServerGameMapChanged(packet *ClientChangeGameMap) *ServerGameMapChanged {
	return &ServerGameMapChanged{
		Packet:               Packet{Id: PacketIdServerGameMapChanged},
		MD5:                  packet.MD5,
		AlternativeMD5:       packet.AlternativeMD5,
		MapId:                packet.MapId,
		MapsetId:             packet.MapsetId,
		Name:                 packet.Name,
		GameMode:             packet.Mode,
		DifficultyRating:     packet.DifficultyRating,
		DifficultyRattingAll: packet.DifficultyRatingAll,
		JudgementCount:       packet.JudgementCount,
	}
}
