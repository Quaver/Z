package packets

type ClientFriendship struct {
	Packet
	UserId int               `json:"u"`
	Action FriendsListAction `json:"a"`
}

type FriendsListAction int

const (
	FriendsListActionAdd FriendsListAction = iota
	FriendsListActionRemove
)
