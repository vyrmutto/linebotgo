package auth

import "time"

// Token holds an active LINE session.
type Token struct {
	AuthToken    string
	RefreshToken string
	ExpiresAt    time.Time
	DeviceID     string
	Certificate  string
}

// IsExpired returns true if the token is past its expiry time.
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
