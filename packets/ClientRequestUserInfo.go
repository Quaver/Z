package packets

type ClientRequestUserInfo struct {
	Packet
	UserIds []int `json:"uids"`
}
