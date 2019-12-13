package common

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/getsentry/raven-go"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
	"gopkg.in/redis.v5"
)

// RavenClient is the raven client to which report errors happening.
// If nil, errors will just be fmt.Println'd
var RavenClient *raven.Client

// MethodData is a struct containing the data passed over to an API method.
type MethodData struct {
	User  Token
	DB    *sqlx.DB
	Doggo *statsd.Client
	R     *redis.Client
	Ctx   *fasthttp.RequestCtx
}

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func (md MethodData) ClientIP() string {
	clientIP := strings.TrimSpace(string(md.Ctx.Request.Header.Peek("X-Real-Ip")))
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = string(md.Ctx.Request.Header.Peek("X-Forwarded-For"))
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	return md.Ctx.RemoteIP().String()
}

// Err logs an error. If RavenClient is set, it will use the client to report
// the error to sentry, otherwise it will just write the error to stdout.
func (md MethodData) Err(err error) {
	user := &raven.User{
		ID:       strconv.Itoa(md.User.UserID),
		Username: md.User.Value,
		IP:       md.Ctx.RemoteAddr().String(),
	}
	// Generate tags for error
	tags := map[string]string{
		"endpoint": string(md.Ctx.RequestURI()),
		"token":    md.User.Value,
	}
	_err(err, tags, user, md.Ctx)
}

// Err for peppy API calls
func Err(c *fasthttp.RequestCtx, err error) {
	// Generate tags for error
	tags := map[string]string{
		"endpoint": string(c.RequestURI()),
	}

	_err(err, tags, nil, c)
}

// WSErr is the error function for errors happening in the websockets.
func WSErr(err error) {
	_err(err, map[string]string{
		"endpoint": "/api/v1/ws",
	}, nil, nil)
}

// GenericError is just an error. Can't make a good description.
func GenericError(err error) {
	_err(err, nil, nil, nil)
}

func _err(err error, tags map[string]string, user *raven.User, c *fasthttp.RequestCtx) {
	if RavenClient == nil {
		fmt.Println("ERROR!!!!")
		fmt.Println(err)
		return
	}

	// Create stacktrace
	st := raven.NewStacktrace(0, 3, []string{"zxq.co/ripple", "git.zxq.co/ripple"})

	ifaces := []raven.Interface{st, generateRavenHTTP(c)}
	if user != nil {
		ifaces = append(ifaces, user)
	}

	RavenClient.CaptureError(
		err,
		tags,
		ifaces...,
	)
}

func generateRavenHTTP(ctx *fasthttp.RequestCtx) *raven.Http {
	if ctx == nil {
		return nil
	}

	// build uri
	uri := ctx.URI()
	// safe to use b2s because a new string gets allocated eventually for
	// concatenation
	sURI := b2s(uri.Scheme()) + "://" + b2s(uri.Host()) + b2s(uri.Path())

	// build header map
	// using ctx.Request.Header.Len would mean calling .VisitAll two times
	// which can be quite expensive since it means iterating over all the
	// headers, so we give a rough estimate of the number of headers we expect
	// to have
	m := make(map[string]string, 16)
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		// not using b2s because we mustn't keep references to the underlying
		// k and v
		m[string(k)] = string(v)
	})

	return &raven.Http{
		URL: sURI,
		// Not using b2s because raven sending is concurrent and may happen
		// AFTER the request, meaning that values could potentially be replaced
		// by new ones.
		Method:  string(ctx.Method()),
		Query:   string(uri.QueryString()),
		Cookies: string(ctx.Request.Header.Peek("Cookie")),
		Headers: m,
	}
}

// ID retrieves the Token's owner user ID.
func (md MethodData) ID() int {
	return md.User.UserID
}

// Query is shorthand for md.C.Query.
func (md MethodData) Query(q string) string {
	return b2s(md.Ctx.QueryArgs().Peek(q))
}

// HasQuery returns true if the parameter is encountered in the querystring.
// It returns true even if the parameter is "" (the case of ?param&etc=etc)
func (md MethodData) HasQuery(q string) bool {
	return md.Ctx.QueryArgs().Has(q)
}

// Unmarshal unmarshals a request's JSON body into an interface.
func (md MethodData) Unmarshal(into interface{}) error {
	return json.Unmarshal(md.Ctx.PostBody(), into)
}

// IsBearer tells whether the current token is a Bearer (oauth) token.
func (md MethodData) IsBearer() bool {
	return md.User.ID == -1
}
