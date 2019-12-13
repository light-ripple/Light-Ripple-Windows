package websockets

import (
	"encoding/json"

	"zxq.co/ripple/rippleapi/common"
)

// SetRestrictedVisibility sets whether the information of restricted users
// can be seen.
func SetRestrictedVisibility(c *conn, message incomingMessage) {
	var visibility bool

	err := json.Unmarshal(message.Data, &visibility)
	if err != nil {
		c.WriteJSON(TypeInvalidMessage, err.Error())
		return
	}

	var userIsManager bool
	if c.User != nil && (c.User.UserPrivileges&uint64(common.AdminPrivilegeManageUsers) > 0) {
		userIsManager = true
	}

	c.Mtx.Lock()
	visibility = visibility && userIsManager
	c.RestrictedVisible = visibility
	c.Mtx.Unlock()

	c.WriteJSON(TypeRestrictedVisibilitySet, visibility)
}
