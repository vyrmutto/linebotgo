package client_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vysina/linebotgo/auth"
	"github.com/vysina/linebotgo/client"
)

func newTestClient(t *testing.T, srvURL string) *client.Client {
	t.Helper()
	tok := &auth.Token{AuthToken: "test-token", ExpiresAt: time.Now().Add(time.Hour)}
	return client.New(client.WithToken(tok), client.WithBaseURL(srvURL))
}

func TestSendText(t *testing.T) {
	var gotBody map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := c.SendText("friend-mid", "Hello!")
	assert.NoError(t, err)
	assert.Equal(t, "friend-mid", gotBody["to"])
	assert.Equal(t, "Hello!", gotBody["text"])
}

func TestOnMessage(t *testing.T) {
	c := client.New()
	var received *client.Message
	c.OnMessage(func(msg *client.Message) { received = msg })
	// OnMessage handler registration should not panic
	assert.NotNil(t, c)
	_ = received
}
