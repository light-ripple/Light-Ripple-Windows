package cieca

import (
	"testing"
)

func TestCSRF(t *testing.T) {
	c := NewCSRF()
	tok, err := c.Generate(1009)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := c.Validate(1009, tok)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("ok is false")
	}
}
