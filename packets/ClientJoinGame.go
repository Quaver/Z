package packets

type ClientJoinGame struct {
	Packet
	GameId   string `json:"gid"`
	Password string `json:"p"`
}
