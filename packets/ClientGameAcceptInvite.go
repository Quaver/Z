package packets

type ClientGameAcceptInvite struct {
	Packet
	MatchId string `json:"m"`
}
