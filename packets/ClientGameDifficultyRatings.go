package packets

type ClientGameDifficultyRatings struct {
	Packet
	Difficulties []float64 `json:"d"`
}
