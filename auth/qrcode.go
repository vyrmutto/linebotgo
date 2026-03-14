package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vysina/linebotgo/transport"
)

// QRLogin manages the QR code login flow.
type QRLogin struct {
	client    *transport.HTTPClient
	createURL string
	pollURL   string
}

// NewQRLogin creates a QRLogin. URLs are injectable for testing.
func NewQRLogin(c *transport.HTTPClient, createURL, pollURL string) *QRLogin {
	return &QRLogin{client: c, createURL: createURL, pollURL: pollURL}
}

// Login fetches a QR code URL, calls displayFn with it, then polls until scanned.
// Timeout: 2 minutes.
func (q *QRLogin) Login(displayFn func(qrURL string)) (*Token, error) {
	req, err := http.NewRequest("POST", q.createURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var session struct {
		Session string `json:"session"`
		URL     string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, err
	}

	displayFn(session.URL)

	deadline := time.Now().Add(2 * time.Minute)
	for time.Now().Before(deadline) {
		req, err := http.NewRequest("GET", q.pollURL+"?session="+session.Session, nil)
		if err != nil {
			return nil, err
		}
		resp, err := q.client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusAccepted {
			resp.Body.Close()
			time.Sleep(100 * time.Millisecond) // short sleep in tests
			continue
		}

		var result struct {
			AuthToken    string `json:"authToken"`
			RefreshToken string `json:"refreshToken"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if result.AuthToken != "" {
			return &Token{
				AuthToken:    result.AuthToken,
				RefreshToken: result.RefreshToken,
				ExpiresAt:    time.Now().Add(24 * time.Hour),
			}, nil
		}
	}
	return nil, fmt.Errorf("QR login timed out after 2 minutes")
}
