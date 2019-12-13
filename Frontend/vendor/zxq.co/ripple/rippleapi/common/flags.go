package common

// These are the flags an user can have. Mostly settings or things like whether
// the user has verified their email address.
const (
	FlagEmailVerified = 1 << iota
	FlagCountry2FA
)
