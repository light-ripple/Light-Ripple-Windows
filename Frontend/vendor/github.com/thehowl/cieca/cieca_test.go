package cieca_test

import (
	"testing"
	"time"

	"github.com/thehowl/cieca"
)

func TestSetGet(t *testing.T) {
	x := []byte("carroponte")
	s := &cieca.DataStore{}
	s.Set("test", x)
	if string(s.Get("test")) != "carroponte" {
		t.Fatal("test != carroponte", string(s.Get("test")))
	}
	s.Clean()
}
func TestClean(t *testing.T) {
	s := &cieca.DataStore{}
	s.Set("test", []byte("a"))
	s.Clean()
	if _, ex := s.GetWithExist("test"); ex {
		t.Fatal("value exists even after Clean!")
	}
}
func TestExpire(t *testing.T) {
	s := &cieca.DataStore{}
	defer s.Clean()
	s.SetWithExpiration("zxcvbn", []byte("why?"), time.Nanosecond*100)
	if s.Get("zxcvbn") == nil {
		t.Fatal("Early expiration?")
	}
	if s.Expiration("zxcvbn") == nil {
		t.Fatal("key's expiration is nil")
	}
	time.Sleep(time.Nanosecond * 5000)
	if s.Get("zxcvbn") != nil {
		t.Fatal("Late expiration?")
	}
}
func TestOverwrite(t *testing.T) {
	s := &cieca.DataStore{}
	defer s.Clean()
	s.Set("meme", []byte("1451"))
	s.Set("meme", []byte("1337"))
	if string(s.Get("meme")) != "1337" {
		t.Fatal("No overwrite?")
	}
}
func TestOverwriteWithExpiration(t *testing.T) {
	s := &cieca.DataStore{}
	defer s.Clean()
	s.SetWithExpiration("carroponte", []byte("19689168196"), time.Second)
	s.Delete("carroponte")
	if s.Get("carroponte") != nil {
		t.Fatal("carroponte ain't nil")
	}
}

// just for coverage
func TestCleanOnNil(t *testing.T) {
	var s *cieca.DataStore
	s.Clean()
}
func TestExpirationWhenTheresNoExpiration(t *testing.T) {
	s := &cieca.DataStore{}
	defer s.Clean()
	s.Set("meme", []byte("x"))
	if s.Expiration("NotExisting") != nil || s.Expiration("meme") != nil {
		t.Fatal("WHAT EXPIRATION???")
	}
}

func BenchmarkSetGetDelete(b *testing.B) {
	s := &cieca.DataStore{}
	s.Get("NotExisting")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set("test", []byte("x"))
		s.Get("test")
		s.Delete("test")
	}
}
