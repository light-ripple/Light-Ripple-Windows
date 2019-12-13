package v1

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"zxq.co/ripple/rippleapi/common"
)

// Achievement represents an achievement in the database.
type Achievement struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// LoadAchievementsEvery reloads the achievements in the database every given
// amount of time.
func LoadAchievementsEvery(db *sqlx.DB, d time.Duration) {
	for {
		achievs = nil
		err := db.Select(&achievs,
			"SELECT id, name, description, icon FROM achievements ORDER BY id ASC")
		if err != nil {
			fmt.Println("LoadAchievements error", err)
			common.GenericError(err)
		}
		time.Sleep(d)
	}
}

var achievs []Achievement

type userAchievement struct {
	Achievement
	Achieved bool `json:"achieved"`
}

type userAchievementsResponse struct {
	common.ResponseBase
	Achievements []userAchievement `json:"achievements"`
}

// UserAchievementsGET handles requests for retrieving the achievements of a
// given user.
func UserAchievementsGET(md common.MethodData) common.CodeMessager {
	shouldRet, whereClause, param := whereClauseUser(md, "users")
	if shouldRet != nil {
		return *shouldRet
	}
	var ids []int
	err := md.DB.Select(&ids, `SELECT ua.achievement_id FROM users_achievements ua
INNER JOIN users ON users.id = ua.user_id
WHERE `+whereClause+` ORDER BY ua.achievement_id ASC`, param)
	switch {
	case err == sql.ErrNoRows:
		return common.SimpleResponse(404, "No such user!")
	case err != nil:
		md.Err(err)
		return Err500
	}
	all := md.HasQuery("all")
	resp := userAchievementsResponse{Achievements: make([]userAchievement, 0, len(achievs))}
	for _, ach := range achievs {
		achieved := inInt(ach.ID, ids)
		if all || achieved {
			resp.Achievements = append(resp.Achievements, userAchievement{ach, achieved})
		}
	}
	resp.Code = 200
	return resp
}

func inInt(i int, js []int) bool {
	for _, j := range js {
		if i == j {
			return true
		}
	}
	return false
}
