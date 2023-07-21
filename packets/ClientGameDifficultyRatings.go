package packets

type ClientGameDifficultyRatings struct {
	Packet
	Md5            string    `json:"md5"`
	AlternativeMd5 string    `json:"amd5"`
	Difficulties   []float64 `json:"d"`
}
