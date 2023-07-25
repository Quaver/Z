package packets

import "example.com/Quaver/Z/objects"

type ClientStatus map[int]*objects.ClientStatus

type ServerUserStatus struct {
	Packet
	Statuses ClientStatus `json:"st"`
}

func NewServerUserStatus(userStatuses ClientStatus) *ServerUserStatus {
	return &ServerUserStatus{
		Packet:   Packet{Id: PacketIdServerUserStatus},
		Statuses: userStatuses,
	}
}

func NewServerUserStatusSingle(userId int, status *objects.ClientStatus) *ServerUserStatus {
	return NewServerUserStatus(ClientStatus{
		userId: status,
	})
}
