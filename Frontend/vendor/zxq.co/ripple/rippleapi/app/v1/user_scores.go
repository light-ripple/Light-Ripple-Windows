package v1

import (
	"fmt"
	"strings"

	"gopkg.in/thehowl/go-osuapi.v1"
	"zxq.co/ripple/rippleapi/common"
	"zxq.co/x/getrank"
)

type userScore struct {
	Score
	Beatmap beatmap `json:"beatmap"`
}

type userScoresResponse struct {
	common.ResponseBase
	Scores []userScore `json:"scores"`
}

const userScoreSelectBase = `
SELECT
	scores.id, scores.beatmap_md5, scores.score,
	scores.max_combo, scores.full_combo, scores.mods,
	scores.300_count, scores.100_count, scores.50_count,
	scores.gekis_count, scores.katus_count, scores.misses_count,
	scores.time, scores.play_mode, scores.accuracy, scores.pp,
	scores.completed,

	beatmaps.beatmap_id, beatmaps.beatmapset_id, beatmaps.beatmap_md5,
	beatmaps.song_name, beatmaps.ar, beatmaps.od, beatmaps.difficulty_std,
	beatmaps.difficulty_taiko, beatmaps.difficulty_ctb, beatmaps.difficulty_mania,
	beatmaps.max_combo, beatmaps.hit_length, beatmaps.ranked,
	beatmaps.ranked_status_freezed, beatmaps.latest_update
FROM scores
INNER JOIN beatmaps ON beatmaps.beatmap_md5 = scores.beatmap_md5
INNER JOIN users ON users.id = scores.userid
`

// UserScoresBestGET retrieves the best scores of an user, sorted by PP if
// mode is standard and sorted by ranked score otherwise.
func UserScoresBestGET(md common.MethodData) common.CodeMessager {
	cm, wc, param := whereClauseUser(md, "users")
	if cm != nil {
		return *cm
	}
	mc := genModeClause(md)
	// For all modes that have PP, we leave out 0 PP scores.
	if getMode(md.Query("mode")) != "ctb" {
		mc += " AND scores.pp > 0"
	}
	return scoresPuts(md, fmt.Sprintf(
		`WHERE
			scores.completed = '3'
			AND %s
			%s
			AND `+md.User.OnlyUserPublic(true)+`
		ORDER BY scores.pp DESC, scores.score DESC %s`,
		wc, mc, common.Paginate(md.Query("p"), md.Query("l"), 100),
	), param)
}

// UserScoresRecentGET retrieves an user's latest scores.
func UserScoresRecentGET(md common.MethodData) common.CodeMessager {
	cm, wc, param := whereClauseUser(md, "users")
	if cm != nil {
		return *cm
	}
	return scoresPuts(md, fmt.Sprintf(
		`WHERE
			%s
			%s
			AND `+md.User.OnlyUserPublic(true)+`
		ORDER BY scores.id DESC %s`,
		wc, genModeClause(md), common.Paginate(md.Query("p"), md.Query("l"), 100),
	), param)
}

func scoresPuts(md common.MethodData, whereClause string, params ...interface{}) common.CodeMessager {
	rows, err := md.DB.Query(userScoreSelectBase+whereClause, params...)
	if err != nil {
		md.Err(err)
		return Err500
	}
	var scores []userScore
	for rows.Next() {
		var (
			us userScore
			b  beatmap
		)
		err = rows.Scan(
			&us.ID, &us.BeatmapMD5, &us.Score.Score,
			&us.MaxCombo, &us.FullCombo, &us.Mods,
			&us.Count300, &us.Count100, &us.Count50,
			&us.CountGeki, &us.CountKatu, &us.CountMiss,
			&us.Time, &us.PlayMode, &us.Accuracy, &us.PP,
			&us.Completed,

			&b.BeatmapID, &b.BeatmapsetID, &b.BeatmapMD5,
			&b.SongName, &b.AR, &b.OD, &b.Diff2.STD,
			&b.Diff2.Taiko, &b.Diff2.CTB, &b.Diff2.Mania,
			&b.MaxCombo, &b.HitLength, &b.Ranked,
			&b.RankedStatusFrozen, &b.LatestUpdate,
		)
		if err != nil {
			md.Err(err)
			return Err500
		}
		b.Difficulty = b.Diff2.STD
		us.Beatmap = b
		us.Rank = strings.ToUpper(getrank.GetRank(
			osuapi.Mode(us.PlayMode),
			osuapi.Mods(us.Mods),
			us.Accuracy,
			us.Count300,
			us.Count100,
			us.Count50,
			us.CountMiss,
		))
		scores = append(scores, us)
	}
	r := userScoresResponse{}
	r.Code = 200
	r.Scores = scores
	return r
}
