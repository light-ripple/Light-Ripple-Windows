package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
)

// PeppyMethod generates a method for the peppyapi
func PeppyMethod(a func(c *fasthttp.RequestCtx, db *sqlx.DB)) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		doggo.Incr("requests.peppy", nil, 1)

		c.Response.Header.SetContentType("application/json; charset=utf-8")

		// I have no idea how, but I manged to accidentally string the first 4
		// letters of the alphabet into a single function call.
		a(c, db)
	}
}
