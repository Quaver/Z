package packets

type ClientStartSpectatingPlayer struct {
	Packet
	UserId int `json:"u"`
}
