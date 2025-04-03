package packets

type ServerGameEnablePreview struct {
	Packet
	Enabled bool `json:"h"`
}

func NewServerGameEnablePreview(enabled bool) *ServerGameEnablePreview {
	return &ServerGameEnablePreview{
		Packet:  Packet{Id: PacketIdServerGameEnablePreviewChanged},
		Enabled: enabled,
	}
}
