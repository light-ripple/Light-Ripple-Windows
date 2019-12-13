package v1

import (
	"strings"

	"zxq.co/ripple/rippleapi/common"
	"zxq.co/ripple/semantic-icons-ugc"
)

type donorInfoResponse struct {
	common.ResponseBase
	HasDonor   bool                 `json:"has_donor"`
	Expiration common.UnixTimestamp `json:"expiration"`
}

// UsersSelfDonorInfoGET returns information about the users' donor status
func UsersSelfDonorInfoGET(md common.MethodData) common.CodeMessager {
	var r donorInfoResponse
	var privileges uint64
	err := md.DB.QueryRow("SELECT privileges, donor_expire FROM users WHERE id = ?", md.ID()).
		Scan(&privileges, &r.Expiration)
	if err != nil {
		md.Err(err)
		return Err500
	}
	r.HasDonor = common.UserPrivileges(privileges)&common.UserPrivilegeDonor > 0
	r.Code = 200
	return r
}

type favouriteModeResponse struct {
	common.ResponseBase
	FavouriteMode int `json:"favourite_mode"`
}

// UsersSelfFavouriteModeGET gets the current user's favourite mode
func UsersSelfFavouriteModeGET(md common.MethodData) common.CodeMessager {
	var f favouriteModeResponse
	f.Code = 200
	if md.ID() == 0 {
		return f
	}
	err := md.DB.QueryRow("SELECT users_stats.favourite_mode FROM users_stats WHERE id = ?", md.ID()).
		Scan(&f.FavouriteMode)
	if err != nil {
		md.Err(err)
		return Err500
	}
	return f
}

type userSettingsData struct {
	UsernameAKA   *string `json:"username_aka"`
	FavouriteMode *int    `json:"favourite_mode"`
	CustomBadge   struct {
		singleBadge
		Show *bool `json:"show"`
	} `json:"custom_badge"`
	PlayStyle *int `json:"play_style"`
}

// UsersSelfSettingsPOST allows to modify information about the current user.
func UsersSelfSettingsPOST(md common.MethodData) common.CodeMessager {
	var d userSettingsData
	md.Unmarshal(&d)

	// input sanitisation
	*d.UsernameAKA = common.SanitiseString(*d.UsernameAKA)
	if md.User.UserPrivileges&common.UserPrivilegeDonor > 0 {
		d.CustomBadge.Name = common.SanitiseString(d.CustomBadge.Name)
		d.CustomBadge.Icon = sanitiseIconName(d.CustomBadge.Icon)
	} else {
		d.CustomBadge.singleBadge = singleBadge{}
		d.CustomBadge.Show = nil
	}
	d.FavouriteMode = intPtrIn(0, d.FavouriteMode, 3)

	q := new(common.UpdateQuery).
		Add("s.username_aka", d.UsernameAKA).
		Add("s.favourite_mode", d.FavouriteMode).
		Add("s.custom_badge_name", d.CustomBadge.Name).
		Add("s.custom_badge_icon", d.CustomBadge.Icon).
		Add("s.show_custom_badge", d.CustomBadge.Show).
		Add("s.play_style", d.PlayStyle)
	_, err := md.DB.Exec("UPDATE users u, users_stats s SET "+q.Fields()+" WHERE s.id = u.id AND u.id = ?", append(q.Parameters, md.ID())...)
	if err != nil {
		md.Err(err)
		return Err500
	}
	return UsersSelfSettingsGET(md)
}

func sanitiseIconName(s string) string {
	classes := strings.Split(s, " ")
	n := make([]string, 0, len(classes))
	for _, c := range classes {
		if !in(c, n) && in(c, semanticiconsugc.SaneIcons) {
			n = append(n, c)
		}
	}
	return strings.Join(n, " ")
}

func in(a string, b []string) bool {
	for _, x := range b {
		if x == a {
			return true
		}
	}
	return false
}

type userSettingsResponse struct {
	common.ResponseBase
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Flags    uint   `json:"flags"`
	userSettingsData
}

// UsersSelfSettingsGET allows to get "sensitive" information about the current user.
func UsersSelfSettingsGET(md common.MethodData) common.CodeMessager {
	var r userSettingsResponse
	var ccb bool
	r.Code = 200
	err := md.DB.QueryRow(`
SELECT
	u.id, u.username,
	u.email, s.username_aka, s.favourite_mode,
	s.show_custom_badge, s.custom_badge_icon,
	s.custom_badge_name, s.can_custom_badge,
	s.play_style, u.flags
FROM users u
LEFT JOIN users_stats s ON u.id = s.id
WHERE u.id = ?`, md.ID()).Scan(
		&r.ID, &r.Username,
		&r.Email, &r.UsernameAKA, &r.FavouriteMode,
		&r.CustomBadge.Show, &r.CustomBadge.Icon,
		&r.CustomBadge.Name, &ccb,
		&r.PlayStyle, &r.Flags,
	)
	if err != nil {
		md.Err(err)
		return Err500
	}
	if !ccb {
		r.CustomBadge = struct {
			singleBadge
			Show *bool `json:"show"`
		}{}
	}
	return r
}

func intPtrIn(x int, y *int, z int) *int {
	if y == nil {
		return nil
	}
	if *y > z {
		return nil
	}
	if *y < x {
		return nil
	}
	return y
}
