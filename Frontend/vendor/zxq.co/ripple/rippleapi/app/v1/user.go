// Package v1 implements the first version of the Ripple API.
package v1

import (
	"database/sql"
	"strconv"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
	"zxq.co/ripple/ocl"
	"zxq.co/ripple/rippleapi/common"
)

type userData struct {
	ID             int                  `json:"id"`
	Username       string               `json:"username"`
	UsernameAKA    string               `json:"username_aka"`
	RegisteredOn   common.UnixTimestamp `json:"registered_on"`
	Privileges     uint64               `json:"privileges"`
	LatestActivity common.UnixTimestamp `json:"latest_activity"`
	Country        string               `json:"country"`
}

const userFields = `SELECT users.id, users.username, register_datetime, users.privileges,
	latest_activity, users_stats.username_aka,
	users_stats.country
FROM users
INNER JOIN users_stats
ON users.id=users_stats.id
`

// UsersGET is the API handler for GET /users
func UsersGET(md common.MethodData) common.CodeMessager {
	shouldRet, whereClause, param := whereClauseUser(md, "users")
	if shouldRet != nil {
		return userPutsMulti(md)
	}

	query := userFields + `
WHERE ` + whereClause + ` AND ` + md.User.OnlyUserPublic(true) + `
LIMIT 1`
	return userPutsSingle(md, md.DB.QueryRowx(query, param))
}

type userPutsSingleUserData struct {
	common.ResponseBase
	userData
}

func userPutsSingle(md common.MethodData, row *sqlx.Row) common.CodeMessager {
	var err error
	var user userPutsSingleUserData

	err = row.StructScan(&user.userData)
	switch {
	case err == sql.ErrNoRows:
		return common.SimpleResponse(404, "No such user was found!")
	case err != nil:
		md.Err(err)
		return Err500
	}

	user.Code = 200
	return user
}

type userPutsMultiUserData struct {
	common.ResponseBase
	Users []userData `json:"users"`
}

func userPutsMulti(md common.MethodData) common.CodeMessager {
	pm := md.Ctx.Request.URI().QueryArgs().PeekMulti
	// query composition
	wh := common.
		Where("users.username_safe = ?", common.SafeUsername(md.Query("nname"))).
		Where("users.id = ?", md.Query("iid")).
		Where("users.privileges = ?", md.Query("privileges")).
		Where("users.privileges & ? > 0", md.Query("has_privileges")).
		Where("users.privileges & ? = 0", md.Query("has_not_privileges")).
		Where("users_stats.country = ?", md.Query("country")).
		Where("users_stats.username_aka = ?", md.Query("name_aka")).
		Where("privileges_groups.name = ?", md.Query("privilege_group")).
		In("users.id", pm("ids")...).
		In("users.username_safe", safeUsernameBulk(pm("names"))...).
		In("users_stats.username_aka", pm("names_aka")...).
		In("users_stats.country", pm("countries")...)

	var extraJoin string
	if md.Query("privilege_group") != "" {
		extraJoin = " LEFT JOIN privileges_groups ON users.privileges & privileges_groups.privileges = privileges_groups.privileges "
	}

	query := userFields + extraJoin + wh.ClauseSafe() + " AND " + md.User.OnlyUserPublic(true) +
		" " + common.Sort(md, common.SortConfiguration{
		Allowed: []string{
			"id",
			"username",
			"privileges",
			"donor_expire",
			"latest_activity",
			"silence_end",
		},
		Default: "id ASC",
		Table:   "users",
	}) +
		" " + common.Paginate(md.Query("p"), md.Query("l"), 100)

	// query execution
	rows, err := md.DB.Queryx(query, wh.Params...)
	if err != nil {
		md.Err(err)
		return Err500
	}
	var r userPutsMultiUserData
	for rows.Next() {
		var u userData
		err := rows.StructScan(&u)
		if err != nil {
			md.Err(err)
			continue
		}
		r.Users = append(r.Users, u)
	}
	r.Code = 200
	return r
}

// UserSelfGET is a shortcut for /users/id/self. (/users/self)
func UserSelfGET(md common.MethodData) common.CodeMessager {
	md.Ctx.Request.URI().SetQueryString("id=self")
	return UsersGET(md)
}

func safeUsernameBulk(us [][]byte) [][]byte {
	for _, u := range us {
		for idx, v := range u {
			if v == ' ' {
				u[idx] = '_'
				continue
			}
			u[idx] = byte(unicode.ToLower(rune(v)))
		}
	}
	return us
}

type whatIDResponse struct {
	common.ResponseBase
	ID int `json:"id"`
}

// UserWhatsTheIDGET is an API request that only returns an user's ID.
func UserWhatsTheIDGET(md common.MethodData) common.CodeMessager {
	var (
		r          whatIDResponse
		privileges uint64
	)
	err := md.DB.QueryRow("SELECT id, privileges FROM users WHERE username_safe = ? LIMIT 1", common.SafeUsername(md.Query("name"))).Scan(&r.ID, &privileges)
	if err != nil || ((privileges&uint64(common.UserPrivilegePublic)) == 0 &&
		(md.User.UserPrivileges&common.AdminPrivilegeManageUsers == 0)) {
		return common.SimpleResponse(404, "That user could not be found!")
	}
	r.Code = 200
	return r
}

var modesToReadable = [...]string{
	"std",
	"taiko",
	"ctb",
	"mania",
}

type modeData struct {
	RankedScore            uint64  `json:"ranked_score"`
	TotalScore             uint64  `json:"total_score"`
	PlayCount              int     `json:"playcount"`
	PlayTime               int     `json:"play_time"`
	ReplaysWatched         int     `json:"replays_watched"`
	TotalHits              int     `json:"total_hits"`
	Level                  float64 `json:"level"`
	Accuracy               float64 `json:"accuracy"`
	PP                     int     `json:"pp"`
	GlobalLeaderboardRank  *int    `json:"global_leaderboard_rank"`
	CountryLeaderboardRank *int    `json:"country_leaderboard_rank"`
}
type userFullResponse struct {
	common.ResponseBase
	userData
	STD           modeData              `json:"std"`
	Taiko         modeData              `json:"taiko"`
	CTB           modeData              `json:"ctb"`
	Mania         modeData              `json:"mania"`
	PlayStyle     int                   `json:"play_style"`
	FavouriteMode int                   `json:"favourite_mode"`
	Badges        []singleBadge         `json:"badges"`
	CustomBadge   *singleBadge          `json:"custom_badge"`
	SilenceInfo   silenceInfo           `json:"silence_info"`
	CMNotes       *string               `json:"cm_notes,omitempty"`
	BanDate       *common.UnixTimestamp `json:"ban_date,omitempty"`
	Email         string                `json:"email,omitempty"`
}
type silenceInfo struct {
	Reason string               `json:"reason"`
	End    common.UnixTimestamp `json:"end"`
}

// UserFullGET gets all of an user's information, with one exception: their userpage.
func UserFullGET(md common.MethodData) common.CodeMessager {
	shouldRet, whereClause, param := whereClauseUser(md, "users")
	if shouldRet != nil {
		return *shouldRet
	}

	// Hellest query I've ever done.
	query := `
SELECT
	users.id, users.username, users.register_datetime, users.privileges, users.latest_activity,

	users_stats.username_aka, users_stats.country, users_stats.play_style, users_stats.favourite_mode,

	users_stats.custom_badge_icon, users_stats.custom_badge_name, users_stats.can_custom_badge,
	users_stats.show_custom_badge,

	users_stats.ranked_score_std, users_stats.total_score_std, users_stats.playcount_std,
	users_stats.replays_watched_std, users_stats.total_hits_std,
	users_stats.avg_accuracy_std, users_stats.pp_std, users_stats.playtime_std,

	users_stats.ranked_score_taiko, users_stats.total_score_taiko, users_stats.playcount_taiko,
	users_stats.replays_watched_taiko, users_stats.total_hits_taiko,
	users_stats.avg_accuracy_taiko, users_stats.pp_taiko, users_stats.playtime_taiko,

	users_stats.ranked_score_ctb, users_stats.total_score_ctb, users_stats.playcount_ctb,
	users_stats.replays_watched_ctb, users_stats.total_hits_ctb,
	users_stats.avg_accuracy_ctb, users_stats.pp_ctb, users_stats.playtime_ctb,

	users_stats.ranked_score_mania, users_stats.total_score_mania, users_stats.playcount_mania,
	users_stats.replays_watched_mania, users_stats.total_hits_mania,
	users_stats.avg_accuracy_mania, users_stats.pp_mania, users_stats.playtime_mania,

	users.silence_reason, users.silence_end,
	users.notes, users.ban_datetime, users.email

FROM users
LEFT JOIN users_stats
ON users.id=users_stats.id
WHERE ` + whereClause + ` AND ` + md.User.OnlyUserPublic(true) + `
LIMIT 1
`
	// Fuck.
	r := userFullResponse{}
	var (
		b    singleBadge
		can  bool
		show bool
	)
	err := md.DB.QueryRow(query, param).Scan(
		&r.ID, &r.Username, &r.RegisteredOn, &r.Privileges, &r.LatestActivity,

		&r.UsernameAKA, &r.Country,
		&r.PlayStyle, &r.FavouriteMode,

		&b.Icon, &b.Name, &can, &show,

		&r.STD.RankedScore, &r.STD.TotalScore, &r.STD.PlayCount,
		&r.STD.ReplaysWatched, &r.STD.TotalHits,
		&r.STD.Accuracy, &r.STD.PP, &r.STD.PlayTime,

		&r.Taiko.RankedScore, &r.Taiko.TotalScore, &r.Taiko.PlayCount,
		&r.Taiko.ReplaysWatched, &r.Taiko.TotalHits,
		&r.Taiko.Accuracy, &r.Taiko.PP, &r.Taiko.PlayTime,

		&r.CTB.RankedScore, &r.CTB.TotalScore, &r.CTB.PlayCount,
		&r.CTB.ReplaysWatched, &r.CTB.TotalHits,
		&r.CTB.Accuracy, &r.CTB.PP, &r.CTB.PlayTime,

		&r.Mania.RankedScore, &r.Mania.TotalScore, &r.Mania.PlayCount,
		&r.Mania.ReplaysWatched, &r.Mania.TotalHits,
		&r.Mania.Accuracy, &r.Mania.PP, &r.Mania.PlayTime,

		&r.SilenceInfo.Reason, &r.SilenceInfo.End,
		&r.CMNotes, &r.BanDate, &r.Email,
	)
	switch {
	case err == sql.ErrNoRows:
		return common.SimpleResponse(404, "That user could not be found!")
	case err != nil:
		md.Err(err)
		return Err500
	}

	can = can && show && common.UserPrivileges(r.Privileges)&common.UserPrivilegeDonor > 0
	if can && (b.Name != "" || b.Icon != "") {
		r.CustomBadge = &b
	}

	for modeID, m := range [...]*modeData{&r.STD, &r.Taiko, &r.CTB, &r.Mania} {
		m.Level = ocl.GetLevelPrecise(int64(m.TotalScore))

		if i := leaderboardPosition(md.R, modesToReadable[modeID], r.ID); i != nil {
			m.GlobalLeaderboardRank = i
		}
		if i := countryPosition(md.R, modesToReadable[modeID], r.ID, r.Country); i != nil {
			m.CountryLeaderboardRank = i
		}
	}

	rows, err := md.DB.Query("SELECT b.id, b.name, b.icon FROM user_badges ub "+
		"LEFT JOIN badges b ON ub.badge = b.id WHERE user = ?", r.ID)
	if err != nil {
		md.Err(err)
	}

	for rows.Next() {
		var badge singleBadge
		err := rows.Scan(&badge.ID, &badge.Name, &badge.Icon)
		if err != nil {
			md.Err(err)
			continue
		}
		r.Badges = append(r.Badges, badge)
	}

	if md.User.TokenPrivileges&common.PrivilegeManageUser == 0 {
		r.CMNotes = nil
		r.BanDate = nil
		r.Email = ""
	}

	r.Code = 200
	return r
}

type userpageResponse struct {
	common.ResponseBase
	Userpage *string `json:"userpage"`
}

// UserUserpageGET gets an user's userpage, as in the customisable thing.
func UserUserpageGET(md common.MethodData) common.CodeMessager {
	shouldRet, whereClause, param := whereClauseUser(md, "users_stats")
	if shouldRet != nil {
		return *shouldRet
	}
	var r userpageResponse
	err := md.DB.QueryRow("SELECT userpage_content FROM users_stats WHERE "+whereClause+" LIMIT 1", param).Scan(&r.Userpage)
	switch {
	case err == sql.ErrNoRows:
		return common.SimpleResponse(404, "No such user!")
	case err != nil:
		md.Err(err)
		return Err500
	}
	if r.Userpage == nil {
		r.Userpage = new(string)
	}
	r.Code = 200
	return r
}

// UserSelfUserpagePOST allows to change the current user's userpage.
func UserSelfUserpagePOST(md common.MethodData) common.CodeMessager {
	var d struct {
		Data *string `json:"data"`
	}
	md.Unmarshal(&d)
	if d.Data == nil {
		return ErrMissingField("data")
	}
	cont := common.SanitiseString(*d.Data)
	_, err := md.DB.Exec("UPDATE users_stats SET userpage_content = ? WHERE id = ? LIMIT 1", cont, md.ID())
	if err != nil {
		md.Err(err)
	}
	md.Ctx.URI().SetQueryString("id=self")
	return UserUserpageGET(md)
}

func whereClauseUser(md common.MethodData, tableName string) (*common.CodeMessager, string, interface{}) {
	switch {
	case md.Query("id") == "self":
		return nil, tableName + ".id = ?", md.ID()
	case md.Query("id") != "":
		id, err := strconv.Atoi(md.Query("id"))
		if err != nil {
			a := common.SimpleResponse(400, "please pass a valid user ID")
			return &a, "", nil
		}
		return nil, tableName + ".id = ?", id
	case md.Query("name") != "":
		return nil, tableName + ".username_safe = ?", common.SafeUsername(md.Query("name"))
	}
	a := common.SimpleResponse(400, "you need to pass either querystring parameters name or id")
	return &a, "", nil
}

type userLookupResponse struct {
	common.ResponseBase
	Users []lookupUser `json:"users"`
}
type lookupUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// UserLookupGET does a quick lookup of users beginning with the passed
// querystring value name.
func UserLookupGET(md common.MethodData) common.CodeMessager {
	name := common.SafeUsername(md.Query("name"))
	name = strings.NewReplacer(
		"%", "\\%",
		"_", "\\_",
		"\\", "\\\\",
	).Replace(name)
	if name == "" {
		return common.SimpleResponse(400, "please provide an username to start searching")
	}
	name = "%" + name + "%"

	var email string
	if md.User.TokenPrivileges&common.PrivilegeManageUser != 0 &&
		strings.Contains(md.Query("name"), "@") {
		email = md.Query("name")
	}

	rows, err := md.DB.Query("SELECT users.id, users.username FROM users WHERE "+
		"(username_safe LIKE ? OR email = ?) AND "+
		md.User.OnlyUserPublic(true)+" LIMIT 25", name, email)
	if err != nil {
		md.Err(err)
		return Err500
	}

	var r userLookupResponse
	for rows.Next() {
		var l lookupUser
		err := rows.Scan(&l.ID, &l.Username)
		if err != nil {
			continue // can't be bothered to handle properly
		}
		r.Users = append(r.Users, l)
	}

	r.Code = 200
	return r
}
