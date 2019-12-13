package common

import "fmt"

// Token is an API token.
type Token struct {
	ID              int
	Value           string
	UserID          int
	TokenPrivileges Privileges
	UserPrivileges  UserPrivileges
}

// OnlyUserPublic returns a string containing "(user.privileges & 1 = 1 OR users.id = <userID>)"
// if the user does not have the UserPrivilege AdminManageUsers, and returns "1" otherwise.
func (t Token) OnlyUserPublic(userManagerSeesEverything bool) string {
	if userManagerSeesEverything &&
		t.UserPrivileges&AdminPrivilegeManageUsers == AdminPrivilegeManageUsers {
		return "1"
	}
	// It's safe to use sprintf directly even if it's a query, because UserID is an int.
	return fmt.Sprintf("(users.privileges & 1 = 1 OR users.id = '%d')", t.UserID)
}
