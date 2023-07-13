package packets

type ClientGameTransferHost struct {
	Packet
	UserId int `json:"u"`
}
