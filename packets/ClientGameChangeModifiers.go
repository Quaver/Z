package packets

type ClientGameChangeModifiers struct {
	Packet
	Modifiers        int64   `json:"md"`
	DifficultyRating float64 `json:"d"`
}
