// Package websockets implements functionality related to the API websockets.
package websockets

import (
	"github.com/jmoiron/sqlx"
	"gopkg.in/redis.v5"
)

var (
	red *redis.Client
	db  *sqlx.DB
)

// Start begins websocket functionality
func Start(r *redis.Client, _db *sqlx.DB) error {
	red = r
	db = _db
	go scoreRetriever()
	go matchRetriever()
	return nil
}

func cleanup(connID uint64) {
	scoreSubscriptionsMtx.Lock()
	for idx, el := range scoreSubscriptions {
		if el.Conn.ID == connID {
			scoreSubscriptions[idx] = scoreSubscriptions[len(scoreSubscriptions)-1]
			scoreSubscriptions = scoreSubscriptions[:len(scoreSubscriptions)-1]
			break
		}
	}
	scoreSubscriptionsMtx.Unlock()
	multiSubscriptionsMtx.Lock()
	for idx, el := range multiSubscriptions {
		if el.ID == connID {
			multiSubscriptions[idx] = multiSubscriptions[len(multiSubscriptions)-1]
			multiSubscriptions = multiSubscriptions[:len(multiSubscriptions)-1]
			break
		}
	}
	multiSubscriptionsMtx.Unlock()
}
