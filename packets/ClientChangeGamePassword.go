package packets

type ClientChangeGamePassword struct {
	Packet
	Password string `json:"p"`
}
