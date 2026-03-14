package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vysina/linebotgo/auth"
)

func TestTokenExpiry(t *testing.T) {
	expired := &auth.Token{
		AuthToken: "tok",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	assert.True(t, expired.IsExpired())

	valid := &auth.Token{
		AuthToken: "tok",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	assert.False(t, valid.IsExpired())
}

func TestErrorsAreDefined(t *testing.T) {
	assert.NotNil(t, auth.ErrTokenExpired)
	assert.NotNil(t, auth.ErrInvalidToken)
	assert.NotNil(t, auth.ErrRateLimited)
	assert.NotNil(t, auth.ErrLoginBlocked)
	assert.NotNil(t, auth.ErrLoginFailed)
}
