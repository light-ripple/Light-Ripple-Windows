package v1

import (
	"fmt"
	"time"

	"zxq.co/ripple/rippleapi/common"
)

type rapLogData struct {
	Through string `json:"through"`
	Text    string `json:"text"`
}

type rapLogMessage struct {
	rapLogData
	Author    int       `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

type rapLogResponse struct {
	common.ResponseBase
	rapLogMessage
}

// RAPLogPOST creates a new entry in the RAP logs
func RAPLogPOST(md common.MethodData) common.CodeMessager {
	if md.User.UserPrivileges&common.AdminPrivilegeAccessRAP == 0 {
		return common.SimpleResponse(403, "Got lost, kiddo?")
	}

	var d rapLogData
	if err := md.Unmarshal(&d); err != nil {
		fmt.Println(err)
		return ErrBadJSON
	}

	if d.Text == "" {
		return ErrMissingField("text")
	}
	if d.Through == "" {
		ua := string(md.Ctx.UserAgent())
		if len(ua) > 20 {
			ua = ua[:20] + "…"
		}
		d.Through = "API"
		if ua != "" {
			d.Through += " (" + ua + ")"
		}
	}
	if len(d.Through) > 30 {
		d.Through = d.Through[:30]
	}

	created := time.Now()
	_, err := md.DB.Exec("INSERT INTO rap_logs(userid, text, datetime, through) VALUES (?, ?, ?, ?)",
		md.User.UserID, d.Text, created.Unix(), d.Through)
	if err != nil {
		md.Err(err)
		return Err500
	}

	var resp rapLogResponse
	resp.rapLogData = d
	resp.Author = md.User.UserID
	resp.CreatedAt = created.Truncate(time.Second)
	resp.Code = 200

	return resp
}

func rapLog(md common.MethodData, message string) {
	ua := string(md.Ctx.UserAgent())
	if len(ua) > 20 {
		ua = ua[:20] + "…"
	}
	through := "API"
	if ua != "" {
		through += " (" + ua + ")"
	}

	_, err := md.DB.Exec("INSERT INTO rap_logs(userid, text, datetime, through) VALUES (?, ?, ?, ?)",
		md.User.UserID, message, time.Now().Unix(), through)
	if err != nil {
		md.Err(err)
	}
}
