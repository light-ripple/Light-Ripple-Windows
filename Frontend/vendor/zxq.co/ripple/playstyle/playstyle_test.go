package playstyle

import "testing"

func TestPlayStyle(t *testing.T) {
	ps := PlayStyle((1 << 10) - 1)
	t.Log(ps.String())
}
func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PlayStyle(i).String()
	}
}
