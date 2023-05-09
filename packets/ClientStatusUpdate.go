package packets

import "example.com/Quaver/Z/objects"

type ClientStatusUpdate struct {
	Packet
	Status objects.ClientStatus `json:"st"`
}
