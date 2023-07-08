package packets

type ServerGameDifficultyRangeChanged struct {
	Packet
	Min float32 `json:"mind"`
	Max float32 `json:"maxd"`
}

func NewServerGameDifficultyRangeChanged(min float32, max float32) *ServerGameDifficultyRangeChanged {
	return &ServerGameDifficultyRangeChanged{
		Packet: Packet{Id: PacketIdServerGameDifficultyRangeChanged},
		Min:    min,
		Max:    max,
	}
}
