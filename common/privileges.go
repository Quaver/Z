package common

type Privileges int64

const (
	PrivilegeNormal Privileges = 1 << iota
	PrivilegeKickUsers
	PrivilegeBanUsers
	PrivilegeNotifyUsers
	PrivilegeMuteUsers
	PrivilegeRankMapsets
	PrivilegeViewAdminLogs
	PrivilegeEditUsers
	PrivilegeManageBuilds
	PrivilegeManageAlphaKeys
	PrivilegeManageMapsets
	PrivilegeEnableTournamentMode
	PrivilegeWipeUsers
	PrivilegeEditUsername
	PrivilegeEditFlag
	PrivilegeEditPrivileges
	PrivilegeEditGroups
	PrivilegeEditNotes
	PrivilegeEditAvatar
	PrivilegeViewCrashes
	PrivilegeEditDonate
)

// HasPrivilege Returns if a combination of privileges has a given privilege
func HasPrivilege(privilegeCombo UserGroups, privilege UserGroups) bool {
	return privilegeCombo&privilege != 0
}
