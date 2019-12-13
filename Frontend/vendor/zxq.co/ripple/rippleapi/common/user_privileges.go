package common

import "strings"

// user/admin privileges
const (
	UserPrivilegePublic UserPrivileges = 1 << iota
	UserPrivilegeNormal
	UserPrivilegeDonor
	AdminPrivilegeAccessRAP
	AdminPrivilegeManageUsers
	AdminPrivilegeBanUsers
	AdminPrivilegeSilenceUsers
	AdminPrivilegeWipeUsers
	AdminPrivilegeManageBeatmap
	AdminPrivilegeManageServer
	AdminPrivilegeManageSetting
	AdminPrivilegeManageBetaKey
	AdminPrivilegeManageReport
	AdminPrivilegeManageDocs
	AdminPrivilegeManageBadges
	AdminPrivilegeViewRAPLogs
	AdminPrivilegeManagePrivilege
	AdminPrivilegeSendAlerts
	AdminPrivilegeChatMod
	AdminPrivilegeKickUsers
	UserPrivilegePendingVerification
	UserPrivilegeTournamentStaff
	AdminPrivilegeCaker
)

// UserPrivileges represents a bitwise enum of the privileges of an user.
type UserPrivileges uint64

var userPrivilegeString = [...]string{
	"UserPublic",
	"UserNormal",
	"UserDonor",
	"AdminAccessRAP",
	"AdminManageUsers",
	"AdminBanUsers",
	"AdminSilenceUsers",
	"AdminWipeUsers",
	"AdminManageBeatmap",
	"AdminManageServer",
	"AdminManageSetting",
	"AdminManageBetaKey",
	"AdminManageReport",
	"AdminManageDocs",
	"AdminManageBadges",
	"AdminViewRAPLogs",
	"AdminManagePrivilege",
	"AdminSendAlerts",
	"AdminChatMod",
	"AdminKickUsers",
	"UserPendingVerification",
	"UserTournamentStaff",
}

func (p UserPrivileges) String() string {
	var pvs []string
	for i, v := range userPrivilegeString {
		if uint64(p)&uint64(1<<uint(i)) != 0 {
			pvs = append(pvs, v)
		}
	}
	return strings.Join(pvs, ", ")
}
