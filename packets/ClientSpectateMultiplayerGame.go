package packets

type ClientSpectateMultiplayerGame struct {
	Packet
	GameId   string `json:"gid"`
	Password string `json:"p"`
}
