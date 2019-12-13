package limit

import (
	"fmt"
	"sync"
	"time"
)

// Request is a Request with DefaultLimiter.
func Request(u string, perMinute int) { DefaultLimiter.Request(u, perMinute) }

// NonBlockingRequest is a NonBlockingRequest with DefaultLimiter.
func NonBlockingRequest(u string, perMinute int) bool {
	return DefaultLimiter.NonBlockingRequest(u, perMinute)
}

// DefaultLimiter is the RateLimiter used by the package-level
// Request and NonBlockingRequest.
var DefaultLimiter = &RateLimiter{
	Map:   make(map[string]chan struct{}),
	Mutex: &sync.RWMutex{},
}

// RateLimiter is a simple rate limiter.
type RateLimiter struct {
	Map   map[string]chan struct{}
	Mutex *sync.RWMutex
}

// Request is a simple request. Blocks until it can make the request.
func (s *RateLimiter) Request(u string, perMinute int) {
	s.request(u, perMinute, true)
}

// NonBlockingRequest checks if it can do a request. If it can't, it returns
// false, else it returns true if the request succeded.
func (s *RateLimiter) NonBlockingRequest(u string, perMinute int) bool {
	return s.request(u, perMinute, false)
}

func (s *RateLimiter) request(u string, perMinute int, blocking bool) bool {
	s.check()
	s.Mutex.RLock()
	c, exists := s.Map[u]
	s.Mutex.RUnlock()
	if !exists {
		c = makePrefilledChan(perMinute)
		s.Mutex.Lock()
		// Now that we have exclusive read and write-access, we want to
		// make sure we don't overwrite an existing channel. Otherwise,
		// race conditions and panic happen.
		if cNew, exists := s.Map[u]; exists {
			c = cNew
			s.Mutex.Unlock()
		} else {
			s.Map[u] = c
			s.Mutex.Unlock()
			<-c
			go s.filler(u, perMinute)
		}
	}
	return rcv(c, blocking)
}

// rcv receives from a channel, but if blocking is true it waits til something
// is received and always returns true, otherwise if it can't receive it
// returns false.
func rcv(c chan struct{}, blocking bool) bool {
	if blocking {
		<-c
		return true
	}
	select {
	case <-c:
		return true
	default:
		return false
	}
}

func (s *RateLimiter) filler(el string, perMinute int) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()

	s.Mutex.RLock()
	c := s.Map[el]
	s.Mutex.RUnlock()
	for {
		select {
		case c <- struct{}{}:
			time.Sleep(time.Minute / time.Duration(perMinute))
		default: // c is full
			s.Mutex.Lock()
			close(c)
			delete(s.Map, el)
			s.Mutex.Unlock()
			return
		}
	}
}

// check makes sure the map and the mutex are properly initialised.
func (s *RateLimiter) check() {
	if s.Map == nil {
		s.Map = make(map[string]chan struct{})
	}
	if s.Mutex == nil {
		s.Mutex = new(sync.RWMutex)
	}
}

func makePrefilledChan(l int) chan struct{} {
	c := make(chan struct{}, l)
	for i := 0; i < l; i++ {
		c <- struct{}{}
	}
	return c
}
