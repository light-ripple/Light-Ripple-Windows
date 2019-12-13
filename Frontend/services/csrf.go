package services

// CSRF is a service that avoids Cross Site Request Forgery by giving tokens
// that will then be used to make sure no third party is interfering.
type CSRF interface {
	// Generate generates a new CSRF token for an user.
	Generate(userID int) (string, error)
	// Validate checks the CSRF token is valid, and if it is it deletes it.
	Validate(userID int, key string) (bool, error)
}
