package peppy

import (
	"database/sql"
	_json "encoding/json"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
	"zxq.co/ripple/rippleapi/common"
)

var modes = []string{"std", "taiko", "ctb", "mania"}

var defaultResponse = []struct{}{}

func genmode(m string) string {
	i := genmodei(m)
	return modes[i]
}
func genmodei(m string) int {
	v := common.Int(m)
	if v > 3 || v < 0 {
		v = 0
	}
	return v
}
func rankable(m string) bool {
	x := genmodei(m)
	return x != 2
}

func genUser(c *fasthttp.RequestCtx, db *sqlx.DB) (string, string) {
	var whereClause string
	var p string

	// used in second case of switch
	s, err := strconv.Atoi(query(c, "u"))

	switch {
	// We know for sure that it's an username.
	case query(c, "type") == "string":
		whereClause = "users.username_safe = ?"
		p = common.SafeUsername(query(c, "u"))
	// It could be an user ID, so we look for an user with that username first.
	case err == nil:
		err = db.QueryRow("SELECT id FROM users WHERE id = ? LIMIT 1", s).Scan(&p)
		if err == sql.ErrNoRows {
			// If no user with that userID were found, assume username.
			whereClause = "users.username_safe = ?"
			p = common.SafeUsername(query(c, "u"))
		} else {
			// An user with that userID was found. Thus it's an userID.
			whereClause = "users.id = ?"
		}
	// u contains letters, so it's an username.
	default:
		whereClause = "users.username_safe = ?"
		p = common.SafeUsername(query(c, "u"))
	}
	return whereClause, p
}

func query(c *fasthttp.RequestCtx, s string) string {
	return string(c.QueryArgs().Peek(s))
}

func json(c *fasthttp.RequestCtx, code int, data interface{}) {
	c.SetStatusCode(code)
	d, err := _json.Marshal(data)
	if err != nil {
		panic(err)
	}
	c.Write(d)
}
