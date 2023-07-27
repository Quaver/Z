package packets

type ServerSpectatorReplayFrames struct {
	Packet
	UserId    int                  `json:"u"`
	Status    SpectatorFrameStatus `json:"s"`
	AudioTime float64              `json:"a"`
	Frames    *string              `json:"f"`
}

func NewServerSpectatorReplayFrames(userId int, status SpectatorFrameStatus, audioTime float64, frames *string) *ServerSpectatorReplayFrames {
	return &ServerSpectatorReplayFrames{
		Packet:    Packet{Id: PacketIdServerSpectatorReplayFrames},
		UserId:    userId,
		Status:    status,
		AudioTime: audioTime,
		Frames:    frames,
	}
}
