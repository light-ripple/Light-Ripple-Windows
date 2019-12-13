package v1

import (
	"strings"

	"zxq.co/ripple/rippleapi/common"
)

// Boilerplate errors
var (
	Err500     = common.SimpleResponse(500, "An error occurred. Trying again may work. If it doesn't, yell at this Ripple instance admin and tell them to fix the API.")
	ErrBadJSON = common.SimpleResponse(400, "Your JSON for this request is invalid.")
)

// ErrMissingField generates a response to a request when some fields in the JSON are missing.
func ErrMissingField(missingFields ...string) common.CodeMessager {
	return common.ResponseBase{
		Code:    422, // http://stackoverflow.com/a/10323055/5328069
		Message: "Missing parameters: " + strings.Join(missingFields, ", ") + ".",
	}
}
