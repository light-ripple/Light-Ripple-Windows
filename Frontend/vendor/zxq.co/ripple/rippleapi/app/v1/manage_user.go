package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	redis "gopkg.in/redis.v5"
	"zxq.co/ripple/rippleapi/common"
)

type setAllowedData struct {
	UserID  int `json:"user_id"`
	Allowed int `json:"allowed"`
}

// UserManageSetAllowedPOST allows to set the allowed status of an user.
func UserManageSetAllowedPOST(md common.MethodData) common.CodeMessager {
	var data setAllowedData
	if err := md.Unmarshal(&data); err != nil {
		return ErrBadJSON
	}
	if data.Allowed < 0 || data.Allowed > 2 {
		return common.SimpleResponse(400, "Allowed status must be between 0 and 2")
	}
	var banDatetime int64
	var privsSet string
	if data.Allowed == 0 {
		banDatetime = time.Now().Unix()
		privsSet = "privileges = (privileges & ~3)"
	} else if data.Allowed == 1 {
		banDatetime = 0
		privsSet = "privileges = (privileges | 3)"
	} else if data.Allowed == 2 {
		banDatetime = time.Now().Unix()
		privsSet = "privileges = (privileges | 2) & (privileges & ~1)"
	}
	_, err := md.DB.Exec("UPDATE users SET "+privsSet+", ban_datetime = ? WHERE id = ?", banDatetime, data.UserID)
	if err != nil {
		md.Err(err)
		return Err500
	}
	rapLog(md, fmt.Sprintf("changed UserID:%d's allowed to %d. This was done using the API's terrible ManageSetAllowed.", data.UserID, data.Allowed))
	go fixPrivileges(data.UserID, md.DB)
	query := `
SELECT users.id, users.username, register_datetime, privileges,
	latest_activity, users_stats.username_aka,
	users_stats.country
FROM users
LEFT JOIN users_stats
ON users.id=users_stats.id
WHERE users.id=?
LIMIT 1`
	return userPutsSingle(md, md.DB.QueryRowx(query, data.UserID))
}

type userEditData struct {
	ID            int          `json:"id"`
	Username      *string      `json:"username"`
	UsernameAKA   *string      `json:"username_aka"`
	Privileges    *uint64      `json:"privileges"`
	Country       *string      `json:"country"`
	SilenceInfo   *silenceInfo `json:"silence_info"`
	ResetUserpage bool         `json:"reset_userpage"`
	//ResetAvatar bool `json:"reset_avatar"`
}

var privChangeList = [...]string{
	"banned",
	"locked",
	"restricted",
	"removed all restrictions on",
}

// UserEditPOST allows to edit an user's information.
func UserEditPOST(md common.MethodData) common.CodeMessager {
	var data userEditData
	if err := md.Unmarshal(&data); err != nil {
		fmt.Println(err)
		return ErrBadJSON
	}

	if data.ID == 0 {
		return common.SimpleResponse(404, "That user could not be found")
	}

	var prevUser struct {
		Username   string
		Privileges uint64
	}
	err := md.DB.Get(&prevUser, "SELECT username, privileges FROM users WHERE id = ? LIMIT 1", data.ID)

	switch err {
	case nil: // carry on
	case sql.ErrNoRows:
		return common.SimpleResponse(404, "That user could not be found")
	default:
		md.Err(err)
		return Err500
	}

	const initQuery = "UPDATE users SET\n"
	q := initQuery
	var args []interface{}

	// totally did not realise I had to update some fields in users_stats as well
	// and just copy pasting the above code by prefixing "stats" to every
	// variable
	const statsInitQuery = "UPDATE users_stats SET\n"
	statsQ := statsInitQuery
	var statsArgs []interface{}

	if common.UserPrivileges(prevUser.Privileges)&common.AdminPrivilegeManageUsers != 0 &&
		data.ID != md.User.UserID {
		return common.SimpleResponse(403, "Can't edit that user")
	}

	var isBanned bool
	var isSilenced bool
	if data.Privileges != nil {
		// If we want to modify privileges other than Normal/Public, we need to have
		// the right privilege ourselves and AdminManageUsers won't suffice.
		if (*data.Privileges&^3) != (prevUser.Privileges&^3) &&
			md.User.UserPrivileges&common.AdminPrivilegeManagePrivilege == 0 {
			return common.SimpleResponse(403, "Can't modify user privileges without AdminManagePrivileges")
		}
		q += "privileges = ?,\n"
		args = append(args, *data.Privileges)

		// UserPublic became 0, so banned or restricted
		const uPublic = uint64(common.UserPrivilegePublic)
		if *data.Privileges&uPublic == 0 && prevUser.Privileges&uPublic != 0 {
			q += "ban_datetime = ?,\n"
			args = append(args, time.Now().Unix())
			isBanned = true
		}

		// If we modified other privileges apart from Normal and Public, we use a generic
		// "changed user's privileges". Otherwise, we are more descriptive.
		if *data.Privileges^prevUser.Privileges > 3 {
			rapLog(md, fmt.Sprintf("has changed %s's privileges", prevUser.Username))
		} else {
			rapLog(md, fmt.Sprintf("has %s %s", privChangeList[*data.Privileges&3], prevUser.Username))
		}
	}
	if data.Username != nil {
		if strings.Contains(*data.Username, " ") && strings.Contains(*data.Username, "_") {
			return common.SimpleResponse(400, "Mixed spaces and underscores")
		}
		if usernameAvailable(md, *data.Username, data.ID) {
			return common.SimpleResponse(409, "User with that username exists")
		}
		jsonData, _ := json.Marshal(struct {
			UserID      int    `json:"userID"`
			NewUsername string `json:"newUsername"`
		}{data.ID, *data.Username})
		md.R.Publish("peppy:change_username", string(jsonData))
	}
	if data.UsernameAKA != nil {
		statsQ += "username_aka = ?,\n"
		statsArgs = append(statsArgs, *data.UsernameAKA)
	}
	if data.Country != nil {
		statsQ += "country = ?,\n"
		statsArgs = append(statsArgs, *data.Country)
		rapLog(md, fmt.Sprintf("has changed %s country to %s", prevUser.Username, *data.Country))
		appendToUserNotes(md, "country changed to "+*data.Country, data.ID)
	}
	if data.SilenceInfo != nil && md.User.UserPrivileges&common.AdminPrivilegeSilenceUsers != 0 {
		q += "silence_end = ?, silence_reason = ?,\n"
		args = append(args, time.Time(data.SilenceInfo.End).Unix(), data.SilenceInfo.Reason)
		isSilenced = true
	}
	if data.ResetUserpage {
		statsQ += "userpage_content = '',\n"
	}

	if q != initQuery {
		q = q[:len(q)-2] + " WHERE id = ? LIMIT 1"
		args = append(args, data.ID)
		_, err = md.DB.Exec(q, args...)
		if err != nil {
			md.Err(err)
			return Err500
		}

	}
	if statsQ != statsInitQuery {
		statsQ = statsQ[:len(statsQ)-2] + " WHERE id = ? LIMIT 1"
		statsArgs = append(statsArgs, data.ID)
		_, err = md.DB.Exec(statsQ, statsArgs...)
		if err != nil {
			md.Err(err)
			return Err500
		}
	}

	if isBanned || isSilenced {
		if err := updateBanBancho(md.R, data.ID); err != nil {
			md.Err(err)
			return Err500
		}
	}

	rapLog(md, fmt.Sprintf("has updated user %s", prevUser.Username))

	return userPutsSingle(md, md.DB.QueryRowx(userFields+" WHERE users.id = ? LIMIT 1", data.ID))
}

func updateBanBancho(r *redis.Client, user int) error {
	return r.Publish("peppy:ban", strconv.Itoa(user)).Err()
}

type wipeUserData struct {
	ID    int   `json:"id"`
	Modes []int `json:"modes"`
}

// WipeUserPOST wipes an user's scores.
func WipeUserPOST(md common.MethodData) common.CodeMessager {
	var data wipeUserData
	if err := md.Unmarshal(&data); err != nil {
		return ErrBadJSON
	}
	if data.ID == 0 {
		return ErrMissingField("id")
	}
	if len(data.Modes) == 0 {
		return ErrMissingField("modes")
	}

	var userData struct {
		Username   string
		Privileges uint64
	}
	err := md.DB.Get(&userData, "SELECT username, privileges FROM users WHERE id = ?", data.ID)
	switch err {
	case sql.ErrNoRows:
		return common.SimpleResponse(404, "That user could not be found!")
	case nil: // carry on
	default:
		md.Err(err)
		return Err500
	}

	if common.UserPrivileges(userData.Privileges)&common.AdminPrivilegeManageUsers != 0 {
		return common.SimpleResponse(403, "Can't edit that user")
	}

	tx, err := md.DB.Beginx()
	if err != nil {
		md.Err(err)
		return Err500
	}

	for _, mode := range data.Modes {
		if mode < 0 || mode > 3 {
			continue
		}
		_, err = tx.Exec("INSERT INTO scores_removed SELECT * FROM scores WHERE userid = ? AND play_mode = ?", data.ID, mode)
		if err != nil {
			md.Err(err)
		}
		_, err = tx.Exec("DELETE FROM scores WHERE userid = ? AND play_mode = ?", data.ID, mode)
		if err != nil {
			md.Err(err)
		}
		_, err = tx.Exec(strings.Replace(
			`UPDATE users_stats SET total_score_MODE = 0, ranked_score_MODE = 0, replays_watched_MODE = 0,
			playcount_MODE = 0, avg_accuracy_MODE = 0, total_hits_MODE = 0, level_MODE = 0, pp_MODE = 0
			WHERE id = ?`, "MODE", modesToReadable[mode], -1,
		), data.ID)
		if err != nil {
			md.Err(err)
		}
		_, err = tx.Exec("DELETE FROM users_beatmap_playcount WHERE user_id = ? AND game_mode = ?", data.ID, mode)
		if err != nil {
			md.Err(err)
		}
	}

	if err = tx.Commit(); err != nil {
		md.Err(err)
		return Err500
	}

	rapLog(md, fmt.Sprintf("has wiped %s's account", userData.Username))

	return userPutsSingle(md, md.DB.QueryRowx(userFields+" WHERE users.id = ? LIMIT 1", data.ID))
}

func appendToUserNotes(md common.MethodData, message string, user int) {
	message = "\n[" + time.Now().Format("2006-01-02 15:04:05") + "] API: " + message
	_, err := md.DB.Exec("UPDATE users SET notes = CONCAT(COALESCE(notes, ''), ?) WHERE id = ?",
		message, user)
	if err != nil {
		md.Err(err)
	}
}

func usernameAvailable(md common.MethodData, u string, userID int) (r bool) {
	err := md.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username_safe = ? AND id != ?)", common.SafeUsername(u), userID).Scan(&r)
	if err != nil && err != sql.ErrNoRows {
		md.Err(err)
	}
	return
}
