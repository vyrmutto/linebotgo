package openchat_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vysina/linebotgo/auth"
	"github.com/vysina/linebotgo/client"
	"github.com/vysina/linebotgo/openchat"
)

func newTestClient(t *testing.T, srvURL string) *client.Client {
	t.Helper()
	tok := &auth.Token{AuthToken: "tok", ExpiresAt: time.Now().Add(time.Hour)}
	return client.New(client.WithToken(tok), client.WithBaseURL(srvURL))
}

func TestSearch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "golang", r.URL.Query().Get("q"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"chats": []map[string]interface{}{
				{"id": "chat-1", "name": "Golang Thailand"},
			},
		})
	}))
	defer srv.Close()

	oc := openchat.New(newTestClient(t, srv.URL))
	results, err := oc.Search("golang")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Golang Thailand", results[0].Name)
}

func TestJoin(t *testing.T) {
	var gotBody map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	oc := openchat.New(newTestClient(t, srv.URL))
	err := oc.Join("chat-1")
	assert.NoError(t, err)
	assert.Equal(t, "chat-1", gotBody["chatId"])
}

func TestSendText(t *testing.T) {
	var gotBody map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	oc := openchat.New(newTestClient(t, srv.URL))
	err := oc.SendText("chat-1", "Hello OpenChat!")
	assert.NoError(t, err)
	assert.Equal(t, "Hello OpenChat!", gotBody["text"])
}

func TestGetMembers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "chat-1", r.URL.Query().Get("chatId"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"members": []map[string]interface{}{
				{"mid": "u123", "displayName": "Alice", "role": "ADMIN"},
				{"mid": "u456", "displayName": "Bob", "role": "MEMBER"},
			},
		})
	}))
	defer srv.Close()

	oc := openchat.New(newTestClient(t, srv.URL))
	members, err := oc.GetMembers("chat-1")
	assert.NoError(t, err)
	assert.Len(t, members, 2)
	assert.Equal(t, "Alice", members[0].DisplayName)
	assert.Equal(t, "ADMIN", members[0].Role)
}
