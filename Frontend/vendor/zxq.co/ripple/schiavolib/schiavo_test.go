package schiavo

import "testing"

func TestSend(t *testing.T) {
	err := Bunker.Send("onii-chan be gentle pls >///< '); DROP TABLE users;-- **hello markdown!** Just testing schiavolib ~")
	if err != nil {
		t.Fatal(err)
	}
}

func Example() {
	Bunker.Send("Hello world!")
}
