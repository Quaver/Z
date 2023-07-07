package packets

type ClientChangeGameName struct {
	Packet
	Name string `json:"n"`
}
