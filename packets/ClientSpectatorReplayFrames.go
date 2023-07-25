package packets

type ClientSpectatorReplayFrames struct {
	Packet
	Status    SpectatorFrameStatus `json:"s"`
	AudioTime float64              `json:"a"`
	Frames    *string              `json:"f"` // *string because it can also be nil instead of empty.
}

type SpectatorFrameStatus int

const (
	SpectatorFrameSelectingSong SpectatorFrameStatus = iota
	SpectatorFrameNewSong
	SpectatorFramePlaying
	SpectatorFramePaused
	SpectatorFrameFinishedPlaying
)
