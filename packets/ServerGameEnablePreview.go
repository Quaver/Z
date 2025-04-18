package packets

type ServerGameEnablePreview struct {
	Packet
	Enabled bool `json:"e"`
}

func NewServerGameEnablePreview(enabled bool) *ServerGameEnablePreview {
	return &ServerGameEnablePreview{
		Packet:  Packet{Id: PacketIdServerGameEnablePreviewChanged},
		Enabled: enabled,
	}
}
