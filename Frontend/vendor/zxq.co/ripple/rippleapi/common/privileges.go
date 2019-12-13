package common

import "strings"

// These are the various privileges a token can have.
const (
	PrivilegeRead             = 1 << iota // used to be to fetch public data, such as user information etc. this is deprecated.
	PrivilegeReadConfidential             // (eventual) private messages, reports... of self
	PrivilegeWrite                        // change user information, write into confidential stuff...
	PrivilegeManageBadges                 // can change various users' badges.
	PrivilegeBetaKeys                     // can add, remove, upgrade/downgrade, make public beta keys.
	PrivilegeManageSettings               // maintainance, set registrations, global alerts, bancho settings
	PrivilegeViewUserAdvanced             // can see user email, and perhaps warnings in the future, basically.
	PrivilegeManageUser                   // can change user email, allowed status, userpage, rank, username...
	PrivilegeManageRoles                  // translates as admin, as they can basically assign roles to anyone, even themselves
	PrivilegeManageAPIKeys                // admin permission to manage user permission, not only self permissions. Only ever do this if you completely trust the application, because this essentially means to put the entire ripple database in the hands of a (potentially evil?) application.
	PrivilegeBlog                         // can do pretty much anything to the blog, and the documentation.
	PrivilegeAPIMeta                      // can do /meta API calls. basically means they can restart the API server.
	PrivilegeBeatmap                      // rank/unrank beatmaps. also BAT when implemented
	PrivilegeBancho									      // can log in to bancho and use the chat through the delta ws api
)

// Privileges is a bitwise enum of the privileges of an user's API key.
type Privileges uint64

var privilegeString = [...]string{
	"Read",
	"ReadConfidential",
	"Write",
	"ManageBadges",
	"BetaKeys",
	"ManageSettings",
	"ViewUserAdvanced",
	"ManageUser",
	"ManageRoles",
	"ManageAPIKeys",
	"Blog",
	"APIMeta",
	"Beatmap",
	"Bancho",
}

func (p Privileges) String() string {
	var pvs []string
	for i, v := range privilegeString {
		if uint64(p)&uint64(1<<uint(i)) != 0 {
			pvs = append(pvs, v)
		}
	}
	return strings.Join(pvs, ", ")
}

var privilegeMustBe = [...]UserPrivileges{
	1 << 30, // read is deprecated, and should be given out to no-one.
	UserPrivilegeNormal,
	UserPrivilegeNormal,
	AdminPrivilegeAccessRAP | AdminPrivilegeManageBadges,
	AdminPrivilegeAccessRAP | AdminPrivilegeManageBetaKey,
	AdminPrivilegeAccessRAP | AdminPrivilegeManageSetting,
	AdminPrivilegeAccessRAP,
	AdminPrivilegeAccessRAP | AdminPrivilegeManageUsers | AdminPrivilegeBanUsers,
	AdminPrivilegeAccessRAP | AdminPrivilegeManageUsers | AdminPrivilegeManagePrivilege,
	AdminPrivilegeAccessRAP | AdminPrivilegeManageUsers | AdminPrivilegeManageServer,
	AdminPrivilegeChatMod, // temporary?
	AdminPrivilegeManageServer,
	AdminPrivilegeAccessRAP | AdminPrivilegeManageBeatmap,
	UserPrivilegeNormal,
}

// CanOnly removes any privilege that the user has requested to have, but cannot have due to their rank.
func (p Privileges) CanOnly(userPrivs UserPrivileges) Privileges {
	newPrivilege := 0
	for i, v := range privilegeMustBe {
		wants := p&1 == 1
		can := userPrivs&v == v
		if wants && can {
			newPrivilege |= 1 << uint(i)
		}
		p >>= 1
	}
	return Privileges(newPrivilege)
}

var privilegeMap = map[string]Privileges{
	"read_confidential": PrivilegeReadConfidential,
	"write":             PrivilegeWrite,
	"bancho":						 PrivilegeBancho,
}

// OAuthPrivileges returns the equivalent in Privileges of a space-separated
// list of scopes.
func OAuthPrivileges(scopes string) Privileges {
	var p Privileges
	for _, x := range strings.Split(scopes, " ") {
		p |= privilegeMap[x]
	}
	return p
}
