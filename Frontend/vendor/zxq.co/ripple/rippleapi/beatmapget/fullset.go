package beatmapget

import (
	"database/sql"
	"errors"
	"time"

	"zxq.co/ripple/rippleapi/common"
	"gopkg.in/thehowl/go-osuapi.v1"
)

// Set checks if an update is required for all beatmaps in a set.
func Set(s int) error {
	var (
		lastUpdated common.UnixTimestamp
		ranked      int
	)
	err := DB.QueryRow("SELECT latest_update, ranked FROM beatmaps WHERE beatmapset_id = ? LIMIT 1", s).
		Scan(&lastUpdated, &ranked)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return set(s, lastUpdated, ranked)
}

// ErrBeatmapNotFound is returned by Beatmap if a beatmap could not be found.
var ErrBeatmapNotFound = errors.New("beatmapget: beatmap not found")

// Beatmap check if an update is required for all beatmaps in the set
// containing this beatmap.
func Beatmap(b int) (int, error) {
	var (
		setID       int
		lastUpdated common.UnixTimestamp
		ranked      int
	)
	err := DB.QueryRow("SELECT beatmapset_id, latest_update, ranked FROM beatmaps WHERE beatmap_id = ? LIMIT 1", b).
		Scan(&setID, &lastUpdated, &ranked)
	switch err {
	case nil:
		return setID, set(setID, lastUpdated, ranked)
	case sql.ErrNoRows:
		beatmaps, err := Client.GetBeatmaps(osuapi.GetBeatmapsOpts{
			BeatmapID: b,
		})
		if err != nil {
			return 0, err
		}
		if len(beatmaps) == 0 {
			return 0, ErrBeatmapNotFound
		}
		return beatmaps[0].BeatmapSetID, set(beatmaps[0].BeatmapSetID, common.UnixTimestamp(time.Time{}), 0)
	default:
		return setID, err
	}
}

func set(s int, updated common.UnixTimestamp, ranked int) error {
	expire := Expire
	if ranked == 2 {
		expire *= 6
	}
	if time.Now().Before(time.Time(updated).Add(expire)) {
		return nil
	}
	beatmaps, err := Client.GetBeatmaps(osuapi.GetBeatmapsOpts{
		BeatmapSetID: s,
	})
	if err != nil {
		return err
	}
	for _, beatmap := range beatmaps {
		err := UpdateIfRequired(BeatmapDefiningQuality{
			ID: beatmap.BeatmapID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
