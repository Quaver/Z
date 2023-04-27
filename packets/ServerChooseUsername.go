package packets

type ServerChooseUsername struct {
	Packet
}

func NewServerChooseUsername() *ServerChooseUsername {
	return &ServerChooseUsername{
		Packet{Id: PacketIdServerChooseUsername},
	}
}
