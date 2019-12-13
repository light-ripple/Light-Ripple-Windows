// Package cieca is an in-memory Key-Value datastore for Go. The main purpose
// is to easily cache data in the RAM for a certain amount of time.
package cieca

import (
	"sync"
	"time"
)

type element struct {
	data   []byte
	expire time.Time
	cancel chan struct{}
}

// DataStore is the datastore containing all the data.
type DataStore struct {
	data  map[string]element
	mutex *sync.RWMutex
}

// Get retrieves a value in the datastore. If it is not found, nil is returned.
func (d *DataStore) Get(key string) []byte {
	v, _ := d.GetWithExist(key)
	return v
}

// GetWithExist retrieves the value of a key, and a boolean indicating whether
// the value existed in the store.
func (d *DataStore) GetWithExist(key string) ([]byte, bool) {
	d.setup()
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	el, ex := d.data[key]
	return el.data, ex
}

// Expiration gets the time of expiration of a key.
func (d *DataStore) Expiration(key string) *time.Time {
	d.setup()
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	el, exist := d.data[key]
	if !exist || el.cancel == nil {
		return nil
	}
	t := el.expire
	return &t
}

// Delete removes an element from the datastore.
func (d *DataStore) Delete(key string) {
	d.setup()
	d.mutex.Lock()
	defer d.mutex.Unlock()
	exp := d.data[key].cancel
	if exp != nil && time.Now().Sub(d.data[key].expire) <= 0 {
		exp <- struct{}{}
		close(exp)
	}
	delete(d.data, key)
}

// Set sets a key in the datastore with no expiration.
func (d *DataStore) Set(key string, value []byte) {
	d.SetWithExpiration(key, value, -1)
}

// SetWithExpiration sets a key in the datastore with a certain value.
func (d *DataStore) SetWithExpiration(key string, value []byte, expiration time.Duration) {
	d.setup()
	var c chan struct{}
	if expiration >= 0 {
		c = make(chan struct{})
	}
	el := element{
		data:   value,
		expire: time.Now().Add(expiration),
		cancel: c,
	}
	d.Delete(key)
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.data[key] = el
	if c != nil {
		go d.queueDeletion(expiration, c, key)
	}
}

// Clean clears the datastore. Not so thread-safe. Use with care.
func (d *DataStore) Clean() {
	if d == nil {
		return
	}
	for el := range d.data {
		d.Delete(el)
	}
}

func (d *DataStore) queueDeletion(dur time.Duration, canc <-chan struct{}, key string) {
	select {
	case <-time.NewTimer(dur).C:
		d.Delete(key)
	case <-canc:
		// useless, but explicits what we're doing
		return
	}
}

func (d *DataStore) setup() {
	if d.data == nil {
		d.data = make(map[string]element)
	}
	if d.mutex == nil {
		d.mutex = &sync.RWMutex{}
	}
}
