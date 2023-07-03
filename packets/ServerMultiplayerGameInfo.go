package packets

import "example.com/Quaver/Z/objects"

type ServerMultiplayerGameInfo struct {
	Packet
	Game *objects.MultiplayerGame `json:"m"`
}

func NewServerMultiplayerGameInfo(game *objects.MultiplayerGame) *ServerMultiplayerGameInfo {
	return &ServerMultiplayerGameInfo{
		Packet: Packet{Id: PacketIdServerMultiplayerGameInfo},
		Game:   game,
	}
}
