package packets

type ClientGameChangePlayerModifiers struct {
	Packet
	Modifiers int64 `json:"m"`
}
