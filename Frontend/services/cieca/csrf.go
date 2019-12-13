package cieca

import (
	"strconv"
	"time"

	"github.com/thehowl/cieca"
	"zxq.co/ripple/hanayo/services"
	"zxq.co/x/rs"
)

// NewCSRF creates a new CSRF service as described in the services.CSRF
// interface.
func NewCSRF() services.CSRF {
	return &ciecaCSRF{
		DataStore: new(cieca.DataStore),
	}
}

type ciecaCSRF struct {
	*cieca.DataStore
}

func (c *ciecaCSRF) Generate(u int) (string, error) {
	var s string
	for {
		s = rs.String(10)
		_, e := c.GetWithExist(s)
		if !e {
			break
		}
	}
	c.SetWithExpiration(strconv.Itoa(u)+s, nil, time.Minute*15)
	return s, nil
}

func (c *ciecaCSRF) Validate(u int, token string) (bool, error) {
	_, e := c.GetWithExist(strconv.Itoa(u) + token)
	return e, nil
}
