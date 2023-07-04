package objects

type MultiplayerGamePlayerMods struct {
	Id        int   `json:"uid"`
	Modifiers int64 `json:"m"`
}

func NewMultiplayerGamePlayerMods(userId int) MultiplayerGamePlayerMods {
	return MultiplayerGamePlayerMods{
		Id:        userId,
		Modifiers: 0,
	}
}
