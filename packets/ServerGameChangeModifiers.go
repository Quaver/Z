package packets

type ServerGameChangeModifiers struct {
	Packet
	Modifiers        int64   `json:"md"`
	DifficultyRating float64 `json:"d"`
}

func NewServerGameChangeModifiers(modifiers int64, difficultyRating float64) *ServerGameChangeModifiers {
	return &ServerGameChangeModifiers{
		Packet:           Packet{Id: PacketIdServerGameChangeModifiers},
		Modifiers:        modifiers,
		DifficultyRating: difficultyRating,
	}
}
