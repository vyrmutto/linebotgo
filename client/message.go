package client

import (
	"encoding/json"
	"net/http"

	"github.com/vysina/linebotgo/api"
)

func (c *Client) resolveSendURL() string {
	if c.sendURL != "" {
		return c.sendURL
	}
	if c.baseURL != "" {
		return c.baseURL + "/send"
	}
	return api.EndpointTalk
}

// SendText sends a plain text message to a LINE user or group.
func (c *Client) SendText(to, text string) error {
	return c.PostJSON(c.resolveSendURL(), map[string]interface{}{
		"to":   to,
		"text": text,
		"type": string(api.MessageTypeText),
	})
}

// OnMessage registers a handler called for each incoming message.
func (c *Client) OnMessage(fn func(*Message)) {
	c.msgHandlers = append(c.msgHandlers, fn)
}

// dispatchMessage calls all registered OnMessage handlers.
func (c *Client) dispatchMessage(msg *Message) {
	for _, h := range c.msgHandlers {
		h(msg)
	}
}

// GetContacts returns a list of contact MIDs.
func (c *Client) GetContacts() ([]string, error) {
	url := api.EndpointTalk
	if c.baseURL != "" {
		url = c.baseURL + "/contacts"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Line-Access", c.AuthToken())
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result struct {
		IDs []string `json:"ids"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.IDs, nil
}
