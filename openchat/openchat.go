package openchat

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vysina/linebotgo/api"
	"github.com/vysina/linebotgo/client"
)

// OpenChat provides LINE OpenChat (formerly LINE Square) operations.
type OpenChat struct {
	c *client.Client
}

// New creates an OpenChat accessor from an existing client.
func New(c *client.Client) *OpenChat {
	return &OpenChat{c: c}
}

// resolveURL returns endpoint with base URL override if set (for tests).
func (oc *OpenChat) resolveURL(endpoint string) string {
	base := oc.c.BaseURL()
	if base == "" {
		return endpoint
	}
	// Strip the LINE base URL and prepend test base URL
	path := strings.TrimPrefix(endpoint, api.BaseURL)
	return base + path
}

// Search finds OpenChat rooms matching the query string.
func (oc *OpenChat) Search(query string) ([]api.OpenChatInfo, error) {
	req, err := http.NewRequest("GET", oc.resolveURL(api.EndpointOpenChatSearch), nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", query)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("X-Line-Access", oc.c.AuthToken())

	resp, err := oc.c.HTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Chats []api.OpenChatInfo `json:"chats"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Chats, nil
}

// Join joins an OpenChat room by ID.
func (oc *OpenChat) Join(chatID string) error {
	return oc.c.PostJSON(oc.resolveURL(api.EndpointOpenChatJoin), map[string]string{"chatId": chatID})
}

// Leave leaves an OpenChat room.
func (oc *OpenChat) Leave(chatID string) error {
	return oc.c.PostJSON(oc.resolveURL(api.EndpointOpenChatLeave), map[string]string{"chatId": chatID})
}

// SendText sends a text message to an OpenChat room.
func (oc *OpenChat) SendText(chatID, text string) error {
	return oc.c.PostJSON(oc.resolveURL(api.EndpointOpenChatSend), map[string]string{
		"chatId": chatID,
		"text":   text,
	})
}

// OnMessage registers a handler for incoming messages from a specific OpenChat room.
func (oc *OpenChat) OnMessage(chatID string, fn func(*client.OpenChatMessage)) {
	oc.c.RegisterOpenChatHandler(chatID, fn)
}
