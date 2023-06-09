package packets

type ClientRequestUserStatus struct {
	Packet
	UserIds []int `json:"uids"`
}
