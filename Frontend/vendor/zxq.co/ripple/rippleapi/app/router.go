package app

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
	"zxq.co/ripple/rippleapi/common"
)

type router struct {
	r *fasthttprouter.Router
}

func (r router) Method(path string, f func(md common.MethodData) common.CodeMessager, privilegesNeeded ...int) {
	r.r.GET(path, wrap(Method(f, privilegesNeeded...)))
}
func (r router) POSTMethod(path string, f func(md common.MethodData) common.CodeMessager, privilegesNeeded ...int) {
	r.r.POST(path, wrap(Method(f, privilegesNeeded...)))
}
func (r router) Peppy(path string, a func(c *fasthttp.RequestCtx, db *sqlx.DB)) {
	r.r.GET(path, wrap(PeppyMethod(a)))
}
func (r router) GET(path string, handle fasthttp.RequestHandler) {
	r.r.GET(path, wrap(handle))
}
func (r router) PlainGET(path string, handle fasthttp.RequestHandler) {
	r.r.GET(path, handle)
}

const (
	// \x1b is escape code for ESC
	// <ESC>[<n>m is escape sequence for a certain colour
	// no IP is written out because of the hundreds of possible ways to pass IPs
	// to a request when using a reverse proxy
	// this is partly inspired from gin, though made even more simplistic.
	fmtString = "%s | %15s |\x1b[%sm %3d \x1b[0m %-7s %s\n"
	// a kind of human readable RFC3339
	timeFormat = "2006-01-02 15:04:05"
	// color reference
	// http://misc.flogisoft.com/bash/tip_colors_and_formatting
	colorOk    = "42" // green
	colorError = "41" // red
)

// wrap returns a function that wraps around handle, providing middleware
// functionality to apply to all API calls, which is to say:
// - logging
// - panic recovery (reporting to sentry)
// - gzipping
func wrap(handle fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		start := time.Now()

		doggo.Incr("requests", nil, 1)

		defer func() {
			if rval := recover(); rval != nil {
				var err error
				switch rval := rval.(type) {
				case string:
					err = errors.New(rval)
				case error:
					err = rval
				default:
					err = fmt.Errorf("%v - %#v", rval, rval)
				}
				common.Err(c, err)
				c.SetStatusCode(500)
				c.SetBodyString(`{ "code": 500, "message": "something really bad happened" }`)
			}

			// switch color to colorError if statusCode is in [500;600)
			color := colorOk
			statusCode := c.Response.StatusCode()
			if statusCode >= 500 && statusCode < 600 {
				color = colorError
			}

			if bytes.Contains(c.Request.Header.Peek("Accept-Encoding"), s2b("gzip")) {
				c.Response.Header.Add("Content-Encoding", "gzip")
				c.Response.Header.Add("Vary", "Accept-Encoding")
				b := c.Response.Body()
				c.Response.ResetBody()
				fasthttp.WriteGzip(c.Response.BodyWriter(), b)
			}

			// print stuff
			fmt.Printf(
				fmtString,
				time.Now().Format(timeFormat),
				time.Since(start).String(),
				color,
				statusCode,
				c.Method(),
				c.Path(),
			)
		}()

		handle(c)
	}
}
