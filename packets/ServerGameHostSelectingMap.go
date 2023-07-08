package packets

type ServerGameHostSelectingMap struct {
	Packet
	IsSelecting bool `json:"s"`
}

func NewServerGameHostSelectingMap(isSelecting bool) *ServerGameHostSelectingMap {
	return &ServerGameHostSelectingMap{
		Packet:      Packet{Id: PacketIdServerGameHostSelectingMap},
		IsSelecting: isSelecting,
	}
}
