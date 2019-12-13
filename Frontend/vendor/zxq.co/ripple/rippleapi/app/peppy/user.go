// Package peppy implements the osu! API as defined on the osu-api repository wiki (https://github.com/ppy/osu-api/wiki).
package peppy

import (
	"database/sql"
	"fmt"
	"strconv"

	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/thehowl/go-osuapi"
	"github.com/valyala/fasthttp"
	"gopkg.in/redis.v5"
	"zxq.co/ripple/ocl"
	"zxq.co/ripple/rippleapi/common"
)

// R is a redis client.
var R *redis.Client

// GetUser retrieves general user information.
func GetUser(c *fasthttp.RequestCtx, db *sqlx.DB) {
	if query(c, "u") == "" {
		json(c, 200, defaultResponse)
		return
	}
	var user osuapi.User
	whereClause, p := genUser(c, db)
	whereClause = "WHERE " + whereClause

	mode := genmode(query(c, "m"))

	err := db.QueryRow(fmt.Sprintf(
		`SELECT
			users.id, users.username,
			users_stats.playcount_%s, users_stats.ranked_score_%s, users_stats.total_score_%s,
			users_stats.pp_%s, users_stats.avg_accuracy_%s,
			users_stats.country
		FROM users
		LEFT JOIN users_stats ON users_stats.id = users.id
		%s
		LIMIT 1`,
		mode, mode, mode, mode, mode, whereClause,
	), p).Scan(
		&user.UserID, &user.Username,
		&user.Playcount, &user.RankedScore, &user.TotalScore,
		&user.PP, &user.Accuracy,
		&user.Country,
	)
	if err != nil {
		json(c, 200, defaultResponse)
		if err != sql.ErrNoRows {
			common.Err(c, err)
		}
		return
	}

	user.Rank = int(R.ZRevRank("ripple:leaderboard:"+mode, strconv.Itoa(user.UserID)).Val()) + 1
	user.CountryRank = int(R.ZRevRank("ripple:leaderboard:"+mode+":"+strings.ToLower(user.Country), strconv.Itoa(user.UserID)).Val()) + 1
	user.Level = ocl.GetLevelPrecise(user.TotalScore)

	json(c, 200, []osuapi.User{user})
}
