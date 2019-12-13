package common

import "testing"

func TestSafeUsername(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"noChange", "no_change", "no_change"},
		{"toLower", "Change_Me", "change_me"},
		{"complete", "La_M a m m a_putt na", "la_m_a_m_m_a_putt_na"},
	}
	for _, tt := range tests {
		if got := SafeUsername(tt.arg); got != tt.want {
			t.Errorf("%q. SafeUsername() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
