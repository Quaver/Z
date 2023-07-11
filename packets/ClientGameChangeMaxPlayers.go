package packets

type ClientGameChangeMaxPlayers struct {
	Packet
	Count int `json:"n"`
}
