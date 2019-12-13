package websockets

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"gopkg.in/thehowl/go-osuapi.v1"
	"zxq.co/ripple/rippleapi/app/v1"
	"zxq.co/ripple/rippleapi/common"
	"zxq.co/x/getrank"
)

type subscribeScoresUser struct {
	User  int   `json:"user"`
	Modes []int `json:"modes"`
}

// SubscribeScores subscribes a connection to score updates.
func SubscribeScores(c *conn, message incomingMessage) {
	var ssu []subscribeScoresUser
	err := json.Unmarshal(message.Data, &ssu)
	if err != nil {
		c.WriteJSON(TypeInvalidMessage, err.Error())
		return
	}

	scoreSubscriptionsMtx.Lock()

	var found bool
	for idx, el := range scoreSubscriptions {
		// already exists, change the users
		if el.Conn.ID == c.ID {
			found = true
			scoreSubscriptions[idx].Users = ssu
		}
	}

	// if it was not found, we need to add it
	if !found {
		scoreSubscriptions = append(scoreSubscriptions, scoreSubscription{c, ssu})
	}

	scoreSubscriptionsMtx.Unlock()

	c.WriteJSON(TypeSubscribedToScores, ssu)
}

type scoreSubscription struct {
	Conn  *conn
	Users []subscribeScoresUser
}

var scoreSubscriptions []scoreSubscription
var scoreSubscriptionsMtx = new(sync.RWMutex)

func scoreRetriever() {
	ps, err := red.Subscribe("api:score_submission")
	if err != nil {
		fmt.Println(err)
	}
	for {
		msg, err := ps.ReceiveMessage()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		go handleNewScore(msg.Payload)
	}
}

type scoreUser struct {
	UserID     int    `json:"id"`
	Username   string `json:"username"`
	Privileges uint64 `json:"privileges"`
}

type score struct {
	v1.Score
	scoreUser
}

type scoreJSON struct {
	v1.Score
	UserID int       `json:"user_id"`
	User   scoreUser `json:"user"`
}

func handleNewScore(id string) {
	defer catchPanic()
	var s score
	err := db.Get(&s, `
SELECT
	s.id, s.beatmap_md5, s.score, s.max_combo, s.full_combo, s.mods,
	s.300_count, s.100_count, s.50_count, s.gekis_count, s.katus_count, s.misses_count,
	s.time, s.play_mode, s.accuracy, s.pp, s.completed, s.userid AS user_id,
	u.username, u.privileges
FROM scores s
INNER JOIN users u ON s.userid = u.id
WHERE s.id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Rank = strings.ToUpper(getrank.GetRank(
		osuapi.Mode(s.PlayMode),
		osuapi.Mods(s.Mods),
		s.Accuracy,
		s.Count300,
		s.Count100,
		s.Count50,
		s.CountMiss,
	))

	sj := scoreJSON{
		Score:  s.Score,
		UserID: s.UserID,
		User:   s.scoreUser,
	}

	scoreSubscriptionsMtx.RLock()
	cp := make([]scoreSubscription, len(scoreSubscriptions))
	copy(cp, scoreSubscriptions)
	scoreSubscriptionsMtx.RUnlock()

	for _, el := range cp {
		if len(el.Users) > 0 && !scoreUserValid(el.Users, sj) {
			continue
		}

		if sj.User.Privileges&3 != 3 && !el.Conn.RestrictedVisible {
			continue
		}

		el.Conn.WriteJSON(TypeNewScore, sj)
	}
}

func scoreUserValid(users []subscribeScoresUser, s scoreJSON) bool {
	for _, u := range users {
		if u.User == s.UserID {
			if len(u.Modes) > 0 {
				if !inModes(u.Modes, s.PlayMode) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func inModes(modes []int, i int) bool {
	for _, m := range modes {
		if m == i {
			return true
		}
	}
	return false
}

func catchPanic() {
	r := recover()
	if r != nil {
		switch r := r.(type) {
		case error:
			common.WSErr(r)
		default:
			fmt.Println("PANIC", r)
		}
	}
}
