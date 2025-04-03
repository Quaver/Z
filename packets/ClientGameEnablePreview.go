package packets

type ClientGameEnablePreview struct {
	Packet
	Enabled bool `json:"o"`
}
