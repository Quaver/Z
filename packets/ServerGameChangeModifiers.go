package packets

import "example.com/Quaver/Z/common"

type ServerGameChangeModifiers struct {
	Packet
	Modifiers        common.Mods `json:"md"`
	DifficultyRating float64     `json:"d"`
}

func NewServerGameChangeModifiers(modifiers common.Mods, difficultyRating float64) *ServerGameChangeModifiers {
	return &ServerGameChangeModifiers{
		Packet:           Packet{Id: PacketIdServerGameChangeModifiers},
		Modifiers:        modifiers,
		DifficultyRating: difficultyRating,
	}
}
