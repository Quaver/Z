package packets

type ServerGameLongNotePercent struct {
	Packet
	Min int `json:"mn"`
	Max int `json:"mx"`
}

func NewServerGameLongNotePercent(min int, max int) *ServerGameLongNotePercent {
	return &ServerGameLongNotePercent{
		Packet: Packet{Id: PacketIdServerGameLongNotePercentageChanged},
		Min:    min,
		Max:    max,
	}
}
