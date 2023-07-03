package packets

import "example.com/Quaver/Z/objects"

type ClientCreateGame struct {
	Packet
	Game *objects.MultiplayerGame `json:"g"`
}
