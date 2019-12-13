package peppy

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
	"gopkg.in/thehowl/go-osuapi.v1"
	"zxq.co/ripple/rippleapi/common"
	"zxq.co/x/getrank"
)

// GetUserRecent retrieves an user's recent scores.
func GetUserRecent(c *fasthttp.RequestCtx, db *sqlx.DB) {
	getUserX(c, db, "ORDER BY scores.time DESC", common.InString(1, query(c, "limit"), 50, 10))
}

// GetUserBest retrieves an user's best scores.
func GetUserBest(c *fasthttp.RequestCtx, db *sqlx.DB) {
	var sb string
	if rankable(query(c, "m")) {
		sb = "scores.pp"
	} else {
		sb = "scores.score"
	}
	getUserX(c, db, "AND completed = '3' ORDER BY "+sb+" DESC", common.InString(1, query(c, "limit"), 100, 10))
}

func getUserX(c *fasthttp.RequestCtx, db *sqlx.DB, orderBy string, limit int) {
	whereClause, p := genUser(c, db)
	sqlQuery := fmt.Sprintf(
		`SELECT
			beatmaps.beatmap_id, scores.score, scores.max_combo,
			scores.300_count, scores.100_count, scores.50_count,
			scores.gekis_count, scores.katus_count, scores.misses_count,
			scores.full_combo, scores.mods, users.id, scores.time,
			scores.pp, scores.accuracy
		FROM scores
		LEFT JOIN beatmaps ON beatmaps.beatmap_md5 = scores.beatmap_md5
		LEFT JOIN users ON scores.userid = users.id
		WHERE %s AND scores.play_mode = ? AND users.privileges & 1 > 0
		%s
		LIMIT %d`, whereClause, orderBy, limit,
	)
	scores := make([]osuapi.GUSScore, 0, limit)
	m := genmodei(query(c, "m"))
	rows, err := db.Query(sqlQuery, p, m)
	if err != nil {
		json(c, 200, defaultResponse)
		common.Err(c, err)
		return
	}
	for rows.Next() {
		var (
			curscore osuapi.GUSScore
			rawTime  common.UnixTimestamp
			acc      float64
			fc       bool
			mods     int
			bid      *int
		)
		err := rows.Scan(
			&bid, &curscore.Score.Score, &curscore.MaxCombo,
			&curscore.Count300, &curscore.Count100, &curscore.Count50,
			&curscore.CountGeki, &curscore.CountKatu, &curscore.CountMiss,
			&fc, &mods, &curscore.UserID, &rawTime,
			&curscore.PP, &acc,
		)
		if err != nil {
			json(c, 200, defaultResponse)
			common.Err(c, err)
			return
		}
		if bid == nil {
			curscore.BeatmapID = 0
		} else {
			curscore.BeatmapID = *bid
		}
		curscore.FullCombo = osuapi.OsuBool(fc)
		curscore.Mods = osuapi.Mods(mods)
		curscore.Date = osuapi.MySQLDate(rawTime)
		curscore.Rank = strings.ToUpper(getrank.GetRank(
			osuapi.Mode(m),
			curscore.Mods,
			acc,
			curscore.Count300,
			curscore.Count100,
			curscore.Count50,
			curscore.CountMiss,
		))
		scores = append(scores, curscore)
	}
	json(c, 200, scores)
}
