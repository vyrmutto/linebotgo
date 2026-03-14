package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vysina/linebotgo/auth"
	"github.com/vysina/linebotgo/transport"
)

// Client is the main LINE API client.
type Client struct {
	token         *auth.Token
	http          *transport.HTTPClient
	baseURL       string
	sendURL       string
	msgHandlers   []func(*Message)
	ocMsgHandlers map[string][]func(*OpenChatMessage)
}

// Option configures a Client.
type Option func(*Client)

// WithToken creates a client authenticated with an existing token.
func WithToken(t *auth.Token) Option {
	return func(c *Client) { c.token = t }
}

// WithBaseURL overrides the base URL for all API calls (used in tests).
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithSendURL overrides the send message URL (used in tests).
func WithSendURL(url string) Option {
	return func(c *Client) { c.sendURL = url }
}

// New creates a new Client.
func New(opts ...Option) *Client {
	c := &Client{
		http:          transport.NewHTTPClient(),
		ocMsgHandlers: make(map[string][]func(*OpenChatMessage)),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// AuthToken returns the current auth token string.
func (c *Client) AuthToken() string {
	if c.token == nil {
		return ""
	}
	return c.token.AuthToken
}

// BaseURL returns the base URL override (empty string if none).
func (c *Client) BaseURL() string {
	return c.baseURL
}

// HTTPClient exposes the underlying transport client (used by openchat package).
func (c *Client) HTTPClient() *transport.HTTPClient {
	return c.http
}

// PostJSON sends a POST request with JSON body, injecting auth header.
func (c *Client) PostJSON(url string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("X-Line-Access", c.AuthToken())
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("LINE API error: status %d", resp.StatusCode)
	}
	return nil
}

// RegisterOpenChatHandler registers a handler for OpenChat messages from a specific chat.
func (c *Client) RegisterOpenChatHandler(chatID string, fn func(*OpenChatMessage)) {
	c.ocMsgHandlers[chatID] = append(c.ocMsgHandlers[chatID], fn)
}

// Message is a received LINE message.
type Message struct {
	ID   string
	From string
	To   string
	Text string
	Type string
}

// OpenChatMessage is a received OpenChat message.
type OpenChatMessage struct {
	ChatID string
	From   string
	Text   string
}
