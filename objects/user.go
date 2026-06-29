package objects

import (
	"example.com/Quaver/Z/common"
)

type PacketUser struct {
	Id              int               `json:"id"`
	SteamId         string            `json:"sid"`
	Username        string            `json:"u"`
	UserGroups      common.UserGroups `json:"ug"`
	MuteEndTime     int64             `json:"m"`
	Country         string            `json:"c"`
	ClanId          int               `json:"cid,omitempty"`
	ClanTag         string            `json:"ct,omitempty"`
	ClanAccentColor string            `json:"ca,omitempty"`
}
