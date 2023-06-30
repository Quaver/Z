package objects

type MultiplayerGameFreeMod int

const (
	MultiplayerGameFreeModNone    MultiplayerGameFreeMod = 0
	MultiplayerGameFreeModRegular MultiplayerGameFreeMod = iota << 0
	MultiplayerGameFreeModRate
)
