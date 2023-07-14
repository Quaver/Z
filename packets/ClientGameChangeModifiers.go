package packets

import "example.com/Quaver/Z/common"

type ClientGameChangeModifiers struct {
	Packet
	Modifiers        common.Mods `json:"md"`
	DifficultyRating float64     `json:"d"`
}
