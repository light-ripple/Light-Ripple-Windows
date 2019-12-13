package websockets

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"database/sql"

	"zxq.co/ripple/rippleapi/common"
)

type websocketUser struct {
	ID              int     `json:"id"`
	Username        string  `json:"username"`
	UserPrivileges  uint64  `json:"user_privileges"`
	TokenPrivileges uint64  `json:"token_privileges"`
	ApplicationID   *string `json:"application_id"`
}

type identifyMessage struct {
	Token    string `json:"token"`
	IsBearer bool   `json:"is_bearer"`
}

// Identify sets the identity of the user.
func Identify(c *conn, message incomingMessage) {
	var idMsg identifyMessage
	err := json.Unmarshal(message.Data, &idMsg)
	if err != nil {
		c.WriteJSON(TypeInvalidMessage, err.Error())
		return
	}

	var wsu websocketUser
	if idMsg.IsBearer {
		err = getBearerToken(idMsg.Token, &wsu)
	} else {
		err = db.Get(&wsu, `
SELECT
	t.user as id, t.privileges as token_privileges,
	u.username, u.privileges as user_privileges
FROM tokens t
INNER JOIN users u ON t.user = u.id
WHERE t.token = ?`, fmt.Sprintf("%x", md5.Sum([]byte(idMsg.Token))))
	}

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		c.WriteJSON(TypeNotFound, nil)
		return
	default:
		common.WSErr(err)
		c.WriteJSON(TypeUnexpectedError, nil)
		return
	}

	wsu.TokenPrivileges = uint64(
		common.Privileges(wsu.TokenPrivileges).CanOnly(
			common.UserPrivileges(wsu.UserPrivileges),
		),
	)

	c.Mtx.Lock()
	c.User = &wsu
	c.Mtx.Unlock()

	c.WriteJSON(TypeIdentified, wsu)
}

func getBearerToken(token string, wsu *websocketUser) error {
	var x struct {
		Client string
		Scope  string
		Extra  int
	}
	err := db.Get(&x, "SELECT client, scope, extra FROM osin_access WHERE access_token = ? LIMIT 1", fmt.Sprintf("%x", sha256.Sum256([]byte(token))))
	if err != nil {
		return err
	}

	var userInfo struct {
		Username   string
		Privileges uint64
	}
	err = db.Get(&userInfo, "SELECT username, privileges FROM users WHERE id = ? LIMIT 1", x.Extra)
	if err != nil {
		return err
	}

	wsu.ApplicationID = &x.Client
	wsu.ID = x.Extra
	wsu.Username = userInfo.Username
	wsu.UserPrivileges = userInfo.Privileges
	wsu.TokenPrivileges = uint64(common.OAuthPrivileges(x.Scope))

	return nil
}
