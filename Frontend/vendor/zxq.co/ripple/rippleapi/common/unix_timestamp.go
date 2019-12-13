package common

import (
	"errors"
	"strconv"
	"time"
)

// UnixTimestamp is simply a time.Time, but can be used to convert an
// unix timestamp in the database into a native time.Time.
type UnixTimestamp time.Time

// Scan decodes src into an unix timestamp.
func (u *UnixTimestamp) Scan(src interface{}) error {
	if u == nil {
		return errors.New("rippleapi/common: UnixTimestamp is nil")
	}
	switch src := src.(type) {
	case int64:
		*u = UnixTimestamp(time.Unix(src, 0))
	case float64:
		*u = UnixTimestamp(time.Unix(int64(src), 0))
	case string:
		return u._string(src)
	case []byte:
		return u._string(string(src))
	case nil:
		// Nothing, leave zero value on timestamp
	default:
		return errors.New("rippleapi/common: unhandleable type")
	}
	return nil
}

func (u *UnixTimestamp) _string(s string) error {
	ts, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*u = UnixTimestamp(time.Unix(int64(ts), 0))
	return nil
}

// MarshalJSON -> time.Time.MarshalJSON
func (u UnixTimestamp) MarshalJSON() ([]byte, error) {
	return time.Time(u).MarshalJSON()
}

// UnmarshalJSON -> time.Time.UnmarshalJSON
func (u *UnixTimestamp) UnmarshalJSON(x []byte) error {
	t := new(time.Time)
	err := t.UnmarshalJSON(x)
	*u = UnixTimestamp(*t)
	return err
}
