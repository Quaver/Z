package packets

import "example.com/Quaver/Z/common"

type ClientGameChangePlayerModifiers struct {
	Packet
	Modifiers common.Mods `json:"m"`
}
