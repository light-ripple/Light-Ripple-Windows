package peppy

import (
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/thehowl/go-osuapi"
	"github.com/valyala/fasthttp"
	"zxq.co/ripple/rippleapi/common"
)

// GetBeatmap retrieves general beatmap information.
func GetBeatmap(c *fasthttp.RequestCtx, db *sqlx.DB) {
	var whereClauses []string
	var params []interface{}
	limit := strconv.Itoa(common.InString(1, query(c, "limit"), 500, 500))

	// since value is not stored, silently ignore
	if query(c, "s") != "" {
		whereClauses = append(whereClauses, "beatmaps.beatmapset_id = ?")
		params = append(params, query(c, "s"))
	}
	if query(c, "b") != "" {
		whereClauses = append(whereClauses, "beatmaps.beatmap_id = ?")
		params = append(params, query(c, "b"))
		// b is unique, so change limit to 1
		limit = "1"
	}
	// creator is not stored, silently ignore u and type
	if query(c, "m") != "" {
		m := genmode(query(c, "m"))
		if m == "std" {
			// Since STD beatmaps are converted, all of the diffs must be != 0
			for _, i := range modes {
				whereClauses = append(whereClauses, "beatmaps.difficulty_"+i+" != 0")
			}
		} else {
			whereClauses = append(whereClauses, "beatmaps.difficulty_"+m+" != 0")
			if query(c, "a") == "1" {
				whereClauses = append(whereClauses, "beatmaps.difficulty_std = 0")
			}
		}
	}
	if query(c, "h") != "" {
		whereClauses = append(whereClauses, "beatmaps.beatmap_md5 = ?")
		params = append(params, query(c, "h"))
	}

	where := strings.Join(whereClauses, " AND ")
	if where != "" {
		where = "WHERE " + where
	}

	rows, err := db.Query(`SELECT
	beatmapset_id, beatmap_id, ranked, hit_length,
	song_name, beatmap_md5, ar, od, bpm, playcount,
	passcount, max_combo, difficulty_std, difficulty_taiko, difficulty_ctb, difficulty_mania,
	latest_update

FROM beatmaps `+where+" ORDER BY id DESC LIMIT "+limit,
		params...)
	if err != nil {
		common.Err(c, err)
		json(c, 200, defaultResponse)
		return
	}

	var bms []osuapi.Beatmap
	for rows.Next() {
		var (
			bm              osuapi.Beatmap
			rawRankedStatus int
			rawName         string
			rawLastUpdate   common.UnixTimestamp
			diffs           [4]float64
		)
		err := rows.Scan(
			&bm.BeatmapSetID, &bm.BeatmapID, &rawRankedStatus, &bm.HitLength,
			&rawName, &bm.FileMD5, &bm.ApproachRate, &bm.OverallDifficulty, &bm.BPM, &bm.Playcount,
			&bm.Passcount, &bm.MaxCombo, &diffs[0], &diffs[1], &diffs[2], &diffs[3],
			&rawLastUpdate,
		)
		if err != nil {
			common.Err(c, err)
			continue
		}
		bm.TotalLength = bm.HitLength
		bm.LastUpdate = osuapi.MySQLDate(rawLastUpdate)
		if rawRankedStatus >= 2 {
			bm.ApprovedDate = osuapi.MySQLDate(rawLastUpdate)
		}
		// zero value of ApprovedStatus == osuapi.StatusPending, so /shrug
		bm.Approved = rippleToOsuRankedStatus[rawRankedStatus]
		bm.Artist, bm.Title, bm.DiffName = parseDiffName(rawName)
		for i, diffVal := range diffs {
			if diffVal != 0 {
				bm.Mode = osuapi.Mode(i)
				bm.DifficultyRating = diffVal
				break
			}
		}
		bms = append(bms, bm)
	}

	json(c, 200, bms)
}

var rippleToOsuRankedStatus = map[int]osuapi.ApprovedStatus{
	0: osuapi.StatusPending,
	1: osuapi.StatusWIP, // it means "needs updating", as the one in the db needs to be updated, but whatever
	2: osuapi.StatusRanked,
	3: osuapi.StatusApproved,
	4: osuapi.StatusQualified,
	5: osuapi.StatusLoved,
}

// buggy diffname parser
func parseDiffName(name string) (author string, title string, diffName string) {
	parts := strings.SplitN(name, " - ", 2)
	author = parts[0]
	if len(parts) > 1 {
		title = parts[1]
		if s := strings.Index(title, " ["); s != -1 {
			diffName = title[s+2:]
			if len(diffName) != 0 && diffName[len(diffName)-1] == ']' {
				diffName = diffName[:len(diffName)-1]
			}
			title = title[:s]
		}
	}
	return
}
