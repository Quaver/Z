package packets

type ServerClearSpectateeReplayFrames struct {
	Packet
	UserId int `json:"u"`
}

func NewServerClearSpectateeReplayFrames(userId int) *ServerClearSpectateeReplayFrames {
	return &ServerClearSpectateeReplayFrames{
		Packet: Packet{Id: PacketIdServerStopSpectatePlayer},
		UserId: userId,
	}
}
