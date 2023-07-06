package packets

import "example.com/Quaver/Z/common"

type ClientChangeGameMap struct {
	Packet
	MD5                 string      `json:"md5"`
	AlternativeMD5      string      `json:"amd5"`
	MapId               int         `json:"mid"`
	MapsetId            int         `json:"msid"`
	Name                string      `json:"map"`
	Mode                common.Mode `json:"gm"`
	DifficultyRating    float64     `json:"d"`
	DifficultyRatingAll []float64   `json:"adr"`
	JudgementCount      int         `json:"jc"`
}
