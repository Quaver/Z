package objects

import "example.com/Quaver/Z/common"

type PacketUser struct {
	Id          int               `json:"id"`
	SteamId     string            `json:"sid"`
	Username    string            `json:"u"`
	UserGroups  common.UserGroups `json:"ug"`
	MuteEndTime int64             `json:"m"`
	Country     string            `json:"c"`
}
