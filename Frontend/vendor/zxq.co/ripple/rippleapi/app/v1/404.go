package v1

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
	"zxq.co/ripple/rippleapi/common"
)

type response404 struct {
	common.ResponseBase
	Cats string `json:"cats"`
}

// Handle404 handles requests with no implemented handlers.
func Handle404(c *fasthttp.RequestCtx) {
	c.Response.Header.Add("X-Real-404", "yes")
	data, err := json.MarshalIndent(response404{
		ResponseBase: common.ResponseBase{
			Code: 404,
		},
		Cats: surpriseMe(),
	}, "", "\t")
	if err != nil {
		panic(err)
	}
	c.SetStatusCode(404)
	c.Write(data)
}
