package oauth

import (
	"crypto/sha256"
	"fmt"

	"github.com/RangelReale/osin"
	"github.com/felipeweb/osin-mysql"
)

// storage is a custom type of storage that implements client secrets saved in the dabase as sha256 hashes
type storage struct {
	*mysql.Storage
}

func (s storage) GetClient(id string) (osin.Client, error) {
	cl, err := s.Storage.GetClient(id)
	return client{cl}, err
}

func (s storage) Clone() osin.Storage {
	return s
}

type client struct {
	osin.Client
}

func (c client) ClientSecretMatches(i string) bool {
	a := c.GetSecret() == fmt.Sprintf("%x", sha256.Sum256([]byte(i)))
	//fmt.Println("SES", a)
	return a
}
