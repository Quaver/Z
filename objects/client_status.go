package objects

import "example.com/Quaver/Z/common"

type ClientStatusType int

const (
	ClientStatusInMenus ClientStatusType = iota
	ClientStatusSelecting
	ClientStatusPLaying
	ClientStatusPaused
	ClientStatusWatching
	ClientStatusEditing
	ClientStatusInLobby
	ClientStatusMultiplayer
	ClientStatusListening
)

type ClientStatus struct {
	Status    ClientStatusType `json:"s"`
	MapId     int              `json:"mid"`
	MapMd5    string           `json:"md5"`
	GameMode  common.Mode      `json:"gm"`
	Content   string           `json:"c"`
	Modifiers int64            `json:"mods"`
}
