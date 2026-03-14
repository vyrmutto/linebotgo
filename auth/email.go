package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vysina/linebotgo/transport"
)

// EmailLogin authenticates with LINE using email and password.
// loginURL is the full endpoint — accepts override for testing.
func EmailLogin(c *transport.HTTPClient, loginURL, email, password string) (*Token, error) {
	body := fmt.Sprintf(`{"identifier":{"type":"email","value":%q},"secret":%q}`, email, password)
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, ErrLoginFailed
	case http.StatusTooManyRequests:
		return nil, ErrRateLimited
	case http.StatusOK:
		// continue
	default:
		return nil, fmt.Errorf("LINE login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		AuthToken    string `json:"authToken"`
		RefreshToken string `json:"refreshToken"`
		Certificate  string `json:"certificate"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &Token{
		AuthToken:    result.AuthToken,
		RefreshToken: result.RefreshToken,
		Certificate:  result.Certificate,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}
