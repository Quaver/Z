package packets

type ClientGameAutoHost struct {
	Packet
	Enabled bool `json:"enabled"`
}
