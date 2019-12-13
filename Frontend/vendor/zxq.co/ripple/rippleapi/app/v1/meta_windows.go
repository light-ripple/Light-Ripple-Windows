// +build windows

package v1

import (
	"time"

	"zxq.co/ripple/rippleapi/common"
)

// MetaRestartGET restarts the API with Zero Downtime™.
func MetaRestartGET(md common.MethodData) common.CodeMessager {
	return common.SimpleResponse(200, "brb in your dreams")
}

// MetaKillGET kills the API process. NOTE TO EVERYONE: NEVER. EVER. USE IN PROD.
// Mainly created because I couldn't bother to fire up a terminal, do htop and kill the API each time.
func MetaKillGET(md common.MethodData) common.CodeMessager {
	return common.SimpleResponse(200, "haha")
}

var upSince = time.Now()

type metaUpSinceResponse struct {
	common.ResponseBase
	Code  int   `json:"code"`
	Since int64 `json:"since"`
}

// MetaUpSinceGET retrieves the moment the API application was started.
// Mainly used to get if the API was restarted.
func MetaUpSinceGET(md common.MethodData) common.CodeMessager {
	return metaUpSinceResponse{
		Code:  200,
		Since: int64(upSince.UnixNano()),
	}
}

// MetaUpdateGET updates the API to the latest version, and restarts it.
func MetaUpdateGET(md common.MethodData) common.CodeMessager {
	return common.SimpleResponse(200, "lol u wish")
}
