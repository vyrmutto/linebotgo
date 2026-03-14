package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vysina/linebotgo/auth"
	"github.com/vysina/linebotgo/transport"
)

func TestEmailLogin_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authToken":    "mock-token",
			"refreshToken": "mock-refresh",
			"certificate":  "mock-cert",
		})
	}))
	defer srv.Close()

	c := transport.NewHTTPClient()
	tok, err := auth.EmailLogin(c, srv.URL+"/login", "test@example.com", "password123")

	assert.NoError(t, err)
	assert.Equal(t, "mock-token", tok.AuthToken)
	assert.Equal(t, "mock-refresh", tok.RefreshToken)
	assert.False(t, tok.IsExpired())
}

func TestEmailLogin_BadCredentials(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	c := transport.NewHTTPClient()
	_, err := auth.EmailLogin(c, srv.URL+"/login", "bad@example.com", "wrong")
	assert.ErrorIs(t, err, auth.ErrLoginFailed)
}

func TestEmailLogin_RateLimited(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	c := transport.NewHTTPClient()
	_, err := auth.EmailLogin(c, srv.URL+"/login", "user@example.com", "pass")
	assert.ErrorIs(t, err, auth.ErrRateLimited)
}
