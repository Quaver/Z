package common

type UserGroups int64

const (
	UserGroupNormal = 1 << iota
	UserGroupAdmin
	UserGroupBot
	UserGroupDeveloper
	UserGroupModerator
	UserGroupRankingSupervisor
	UserGroupSwan
	UserGroupContributor
	UserGroupDonator
)

// HasUserGroup Returns if a combination of user groups contains a single group
func HasUserGroup(groupsCombo UserGroups, group UserGroups) bool {
	return groupsCombo&group != 0
}

// HasAnyUserGroup Returns if a combination of group contains any of a slice of groups
func HasAnyUserGroup(groupsCombo UserGroups, groups []UserGroups) bool {
	for _, group := range groups {
		if HasUserGroup(groupsCombo, group) {
			return true
		}
	}

	return false
}

// IsSwan Returns if the user is Swan
func IsSwan(userGroups UserGroups) bool {
	return HasUserGroup(userGroups, UserGroupSwan)
}
