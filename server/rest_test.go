package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vysina/linebotgo/server"
)

// mockSender implements server.Sender for testing.
type mockSender struct {
	lastTo   string
	lastText string
	err      error
}

func (m *mockSender) SendText(to, text string) error {
	m.lastTo = to
	m.lastText = text
	return m.err
}

func TestRESTSendMessage_Success(t *testing.T) {
	ms := &mockSender{}
	srv := server.New(ms, server.Config{})

	body, _ := json.Marshal(map[string]string{"to": "user-mid", "text": "hi"})
	req := httptest.NewRequest(http.MethodPost, "/api/message/send", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "user-mid", ms.lastTo)
	assert.Equal(t, "hi", ms.lastText)
}

func TestRESTSendMessage_MissingFields(t *testing.T) {
	ms := &mockSender{}
	srv := server.New(ms, server.Config{})

	body, _ := json.Marshal(map[string]string{"to": ""})
	req := httptest.NewRequest(http.MethodPost, "/api/message/send", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Handler().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRESTGetContacts(t *testing.T) {
	srv := server.New(&mockSender{}, server.Config{})
	req := httptest.NewRequest(http.MethodGet, "/api/contacts", nil)
	w := httptest.NewRecorder()

	srv.Handler().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)
	assert.Contains(t, result, "contacts")
}
