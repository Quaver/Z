package packets

type JoinGameError int

const (
	JoinGameErrorPassword JoinGameError = iota
	JoinGameErrorFull
	JoinGameErrorMatchNoExists
)

type ServerJoinGameFailed struct {
	Packet
	Error JoinGameError `json:"r"`
}

func NewServerJoinGameFailed(err JoinGameError) *ServerJoinGameFailed {
	return &ServerJoinGameFailed{
		Packet: Packet{Id: PacketIdServerJoinGameFailed},
		Error:  err,
	}
}
