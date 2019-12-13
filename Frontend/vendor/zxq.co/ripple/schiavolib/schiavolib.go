package schiavo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// Channels to which a message can be sent.
const (
	General Channel = "general"
	Bunker  Channel = "bunk"
	ChatLog Channel = "chatlog"
	Staff   Channel = "staff"
	CMs     Channel = "cm"
)

// Channel is just a channel on the discord to which you can send messages.
type Channel string

// SchiavoURL is the base URL for schiavo. Change to var when not hardcoded
var SchiavoURL = ""

// Prefix is a prefix that will be appended to all Schiavo messages if set.
var Prefix = ""

// ForceDo is a meme
var ForceDo bool

var shouldDo = os.Getenv("GIN_MODE") == "release" || os.Getenv("SCHIAVO_LOG") != ""

// Send sends a message to a channel.
func (c Channel) Send(m string) error {
	if !shouldDo && !ForceDo {
		return nil
	}
	if SchiavoURL == "" {
		return nil
	}
	if Prefix != "" {
		m = fmt.Sprintf("**%s** %s", Prefix, m)
	}
	urgay := SchiavoURL + "/" + string(c) + "?message=" + url.QueryEscape(m)
	resp, err := http.Get(urgay)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	if string(body) != "ok" {
		return errors.New("Schiavo response not ok: " + string(body) + "; status code: " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}
