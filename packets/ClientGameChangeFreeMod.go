package packets

import "example.com/Quaver/Z/objects"

type ClientGameFreeModTypeChanged struct {
	Packet
	Type objects.MultiplayerGameFreeMod `json:"t"`
}
