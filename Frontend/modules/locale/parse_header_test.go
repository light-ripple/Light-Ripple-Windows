package locale

import (
	"reflect"
	"testing"
)

func TestParseHeader(t *testing.T) {
	tt := []struct {
		In  string
		Out []string
	}{
		{
			"en",
			[]string{"en"},
		},
		{
			"en-GB",
			[]string{"en_GB"},
		},
		{
			"en-GB;q=0.5,it",
			[]string{"it", "en_GB"},
		},
		{
			"en-GB;q=0.5,it,pl;q=0.2",
			[]string{"it", "en_GB", "pl"},
		},
		{
			"en-GB;q=0.5,pl;q=xd",
			[]string{"pl", "en_GB"},
		},
		{
			"",
			nil,
		},
	}

	for _, el := range tt {
		got := ParseHeader(el.In)
		if !reflect.DeepEqual(got, el.Out) {
			t.Errorf("got %v want %v", got, el.Out)
		}
	}
}
