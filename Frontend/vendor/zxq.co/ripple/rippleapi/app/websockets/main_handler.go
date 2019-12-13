package websockets

import (
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/leavengood/websocket"
)

var stepNumber uint64

func handler(rawConn *websocket.Conn) {
	defer catchPanic()
	defer rawConn.Close()

	step := atomic.AddUint64(&stepNumber, 1)

	// 5 is a security margin in case
	if step == (1<<10 - 5) {
		atomic.StoreUint64(&stepNumber, 0)
	}

	c := &conn{
		rawConn,
		sync.Mutex{},
		step | uint64(time.Now().UnixNano()<<10),
		false,
		nil,
	}

	c.WriteJSON(TypeConnected, nil)

	defer cleanup(c.ID)

	for {
		var i incomingMessage
		err := c.Conn.ReadJSON(&i)
		if _, ok := err.(*websocket.CloseError); ok {
			return
		}
		if err != nil {
			c.WriteJSON(TypeInvalidMessage, err.Error())
			continue
		}
		f, ok := messageHandler[i.Type]
		if !ok {
			c.WriteJSON(TypeInvalidMessage, "invalid message type")
			continue
		}
		if f != nil {
			f(c, i)
		}
	}
}

type conn struct {
	Conn              *websocket.Conn
	Mtx               sync.Mutex
	ID                uint64
	RestrictedVisible bool
	User              *websocketUser
}

func (c *conn) WriteJSON(t string, data interface{}) error {
	c.Mtx.Lock()
	err := c.Conn.WriteJSON(newMessage(t, data))
	c.Mtx.Unlock()
	return err
}

var messageHandler = map[string]func(c *conn, message incomingMessage){
	TypeSubscribeScores:         SubscribeScores,
	TypeSubscribeMultiMatches:   SubscribeMultiMatches,
	TypeSetRestrictedVisibility: SetRestrictedVisibility,
	TypeIdentify:                Identify,
	TypePing:                    pingHandler,
}

// Server Message Types
const (
	TypeConnected                = "connected"
	TypeInvalidMessage           = "invalid_message_type"
	TypeUnexpectedError          = "unexpected_error"
	TypeNotFound                 = "not_found"
	TypeSubscribedToScores       = "subscribed_to_scores"
	TypeNewScore                 = "new_score"
	TypeSubscribedToMultiMatches = "subscribed_mp_complete_match"
	TypeNewMatch                 = "new_completed_match"
	TypeIdentified               = "identified"
	TypeRestrictedVisibilitySet  = "restricted_visibility_set"
	TypePong                     = "pong"
)

// Client Message Types
const (
	TypeSubscribeScores         = "subscribe_scores"
	TypeSubscribeMultiMatches   = "subscribe_mp_complete_match"
	TypeIdentify                = "identify"
	TypeSetRestrictedVisibility = "set_restricted_visibility"
	TypePing                    = "ping"
)

func pingHandler(c *conn, message incomingMessage) {
	c.WriteJSON(TypePong, nil)
}

// Message is the wrapped information for a message sent to the client.
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

func newMessage(t string, data interface{}) Message {
	return Message{
		Type: t,
		Data: data,
	}
}

type incomingMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
