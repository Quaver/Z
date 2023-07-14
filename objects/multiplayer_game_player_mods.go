package objects

import "example.com/Quaver/Z/common"

type MultiplayerGamePlayerMods struct {
	Id        int         `json:"uid"`
	Modifiers common.Mods `json:"m"`
}

func NewMultiplayerGamePlayerMods(userId int) MultiplayerGamePlayerMods {
	return MultiplayerGamePlayerMods{
		Id:        userId,
		Modifiers: 0,
	}
}
