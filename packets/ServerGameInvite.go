package packets

type ServerGameInvite struct {
	Packet
	MatchId string `json:"m"`
	Sender  string `json:"s"`
}

func NewServerGameInvite(matchId string, sender string) *ServerGameInvite {
	return &ServerGameInvite{
		Packet:  Packet{Id: PacketIdServerGameInvite},
		MatchId: matchId,
		Sender:  sender,
	}
}
