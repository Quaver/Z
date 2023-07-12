package packets

type ClientRequestUserStats struct {
	Packet
	Users []int `json:"u"`
}
