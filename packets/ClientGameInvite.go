package packets

type ClientGameInvite struct {
	Packet
	UserId int `json:"u"`
}
