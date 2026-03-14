package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vysina/linebotgo/auth"
	"github.com/vysina/linebotgo/transport"
)

func TestQRLogin_Success(t *testing.T) {
	var callCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		switch n {
		case 1: // create session
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"session": "sess-123",
				"url":     "https://line.me/R/qr?code=abc",
			})
		case 2: // poll — not scanned yet
			w.WriteHeader(http.StatusAccepted)
		default: // poll — scanned
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"authToken":    "qr-token",
				"refreshToken": "qr-refresh",
			})
		}
	}))
	defer srv.Close()

	c := transport.NewHTTPClient()
	qr := auth.NewQRLogin(c, srv.URL+"/create", srv.URL+"/poll")

	var displayedURL string
	tok, err := qr.Login(func(qrURL string) { displayedURL = qrURL })

	assert.NoError(t, err)
	assert.Equal(t, "qr-token", tok.AuthToken)
	assert.Equal(t, "https://line.me/R/qr?code=abc", displayedURL)
	assert.False(t, tok.IsExpired())
}
