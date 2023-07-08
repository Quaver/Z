package packets

type ClientGameHostSelectingMap struct {
	Packet
	IsSelecting bool `json:"s"`
}
