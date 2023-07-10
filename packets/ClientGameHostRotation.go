package packets

type ClientGameHostRotation struct {
	Packet
	Enabled bool `json:"o"`
}
