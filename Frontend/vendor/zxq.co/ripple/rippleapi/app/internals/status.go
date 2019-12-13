// Package internals has methods that suit none of the API packages.
package internals

import "github.com/valyala/fasthttp"

var statusResp = []byte(`{ "status": 1 }`)

// Status is used for checking the API is up by the ripple website, on the status page.
func Status(c *fasthttp.RequestCtx) {
	c.Write(statusResp)
}
