package v1

import (
	"zxq.co/ripple/rippleapi/common"
)

type mostPlayedBeatmap struct {
	Beatmap beatmap `json:"beatmap"`
	PlayCount int	`json:"playcount"`
}

type mostPlayedBeatmapResponse struct {
	common.ResponseBase
	Beatmaps []mostPlayedBeatmap `json:"beatmaps"`
}

func UserMostPlayedGET(md common.MethodData) common.CodeMessager {
	shouldRet, whereClause, param := whereClauseUser(md, "users")
	if shouldRet != nil {
		return *shouldRet
	}
	whereClause += " " + genModeClauseColumn(md, "users_beatmap_playcount.game_mode")
	var q = `SELECT
beatmap_id, beatmapset_id, beatmap_md5,
song_name, ar, od, difficulty_std, difficulty_taiko,
difficulty_ctb, difficulty_mania, max_combo,
hit_length, ranked, ranked_status_freezed,
latest_update, users_beatmap_playcount.playcount 
FROM users_beatmap_playcount LEFT JOIN beatmaps USING(beatmap_id) 
LEFT JOIN users ON users_beatmap_playcount.user_id = users.id 
WHERE ` + whereClause + ` ORDER BY users_beatmap_playcount.playcount DESC ` +
		common.Paginate(md.Query("p"), md.Query("l"), 100)
	rows, err := md.DB.Query(q, param)
	if err != nil {
		md.Err(err)
		return Err500
	}
	var r mostPlayedBeatmapResponse
	for rows.Next() {
		var mpb mostPlayedBeatmap
		err = rows.Scan(
			&mpb.Beatmap.BeatmapID, &mpb.Beatmap.BeatmapsetID, &mpb.Beatmap.BeatmapMD5,
			&mpb.Beatmap.SongName, &mpb.Beatmap.AR, &mpb.Beatmap.OD, &mpb.Beatmap.Diff2.STD, &mpb.Beatmap.Diff2.Taiko,
			&mpb.Beatmap.Diff2.CTB, &mpb.Beatmap.Diff2.Mania, &mpb.Beatmap.MaxCombo,
			&mpb.Beatmap.HitLength, &mpb.Beatmap.Ranked, &mpb.Beatmap.RankedStatusFrozen,
			&mpb.Beatmap.LatestUpdate, &mpb.PlayCount,
		)
		mpb.Beatmap.Difficulty = mpb.Beatmap.Diff2.STD
		if err != nil {
			md.Err(err)
			continue
		}
		r.Beatmaps = append(r.Beatmaps, mpb)
	}
	r.Code = 200
	return r
}
