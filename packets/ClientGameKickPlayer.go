package packets

type ClientGameKickPlayer struct {
	Packet
	UserId int `json:"u"`
}
