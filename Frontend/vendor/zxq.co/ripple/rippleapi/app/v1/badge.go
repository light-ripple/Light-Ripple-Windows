package v1

import (
	"database/sql"

	"zxq.co/ripple/rippleapi/common"
)

type singleBadge struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type multiBadgeData struct {
	common.ResponseBase
	Badges []singleBadge `json:"badges"`
}

// BadgesGET retrieves all the badges on this ripple instance.
func BadgesGET(md common.MethodData) common.CodeMessager {
	var (
		r    multiBadgeData
		rows *sql.Rows
		err  error
	)
	if md.Query("id") != "" {
		rows, err = md.DB.Query("SELECT id, name, icon FROM badges WHERE id = ? LIMIT 1", md.Query("id"))
	} else {
		rows, err = md.DB.Query("SELECT id, name, icon FROM badges " + common.Paginate(md.Query("p"), md.Query("l"), 50))
	}
	if err != nil {
		md.Err(err)
		return Err500
	}
	defer rows.Close()
	for rows.Next() {
		nb := singleBadge{}
		err = rows.Scan(&nb.ID, &nb.Name, &nb.Icon)
		if err != nil {
			md.Err(err)
		}
		r.Badges = append(r.Badges, nb)
	}
	if err := rows.Err(); err != nil {
		md.Err(err)
	}
	r.ResponseBase.Code = 200
	return r
}

type badgeMembersData struct {
	common.ResponseBase
	Members []userData `json:"members"`
}

// BadgeMembersGET retrieves the people who have a certain badge.
func BadgeMembersGET(md common.MethodData) common.CodeMessager {
	i := common.Int(md.Query("id"))
	if i == 0 {
		return ErrMissingField("id")
	}

	var members badgeMembersData

	err := md.DB.Select(&members.Members, `SELECT users.id, users.username, register_datetime, users.privileges,
	latest_activity, users_stats.username_aka,
	users_stats.country
FROM user_badges ub
INNER JOIN users ON users.id = ub.user
INNER JOIN users_stats ON users_stats.id = ub.user
WHERE badge = ?
ORDER BY id ASC `+common.Paginate(md.Query("p"), md.Query("l"), 50), i)

	if err != nil {
		md.Err(err)
		return Err500
	}

	members.Code = 200
	return members
}
