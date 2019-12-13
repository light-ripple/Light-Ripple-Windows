// Package peppy implements the osu! API as defined on the osu-api repository wiki (https://github.com/ppy/osu-api/wiki).
package peppy

import (
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
)

// GetMatch retrieves general match information.
func GetMatch(c *fasthttp.RequestCtx, db *sqlx.DB) {
	json(c, 200, defaultResponse)
}
