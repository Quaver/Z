package packets

type ServerGameTournamentMode struct {
	Packet
	Enabled bool `json:"trn"`
}

func NewServerGameTournamentMode(enabled bool) *ServerGameTournamentMode {
	return &ServerGameTournamentMode{
		Packet:  Packet{Id: PacketIdServerGameTournamentMode},
		Enabled: enabled,
	}
}
