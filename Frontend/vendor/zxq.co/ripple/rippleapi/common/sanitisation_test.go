package common

import "testing"

const pen = "I trattori di palmizio 나는 펜이있다. 私はリンゴを持っています。" +
	"啊! 苹果笔。 у меня есть ручка, Tôi có dứa. අන්නාසි පෑන"

func TestSanitiseString(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			"Normal",
			pen,
			pen,
		},
		{
			"Arabic (rtl)",
			"أناناس",
			"أناناس",
		},
		{
			"Null",
			"A\x00B",
			"AB",
		},
	}
	for _, tt := range tests {
		if got := SanitiseString(tt.arg); got != tt.want {
			t.Errorf("%q. SanitiseString() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func BenchmarkSanitiseString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SanitiseString(pen)
	}
}
