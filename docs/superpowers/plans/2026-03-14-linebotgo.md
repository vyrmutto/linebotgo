# linebotgo Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go library + HTTP/WebSocket server that wraps LINE's private API for sending/receiving messages and interacting with OpenChat without the official LINE Messaging API.

**Architecture:** Layered SDK — `transport/` handles raw HTTP+Thrift to LINE servers, `auth/` manages login and token lifecycle, `api/` holds all LINE API definitions (the layer updated when LINE changes), `client/` exposes the high-level Go SDK, `openchat/` adds OpenChat features, and `server/` wraps the client as a REST+WebSocket server.

**Tech Stack:** Go 1.21+, `net/http` (stdlib), `nhooyr.io/websocket` (WebSocket), `encoding/json` (REST), Apache Thrift (LINE protocol), `github.com/spf13/viper` or plain JSON for config.

**Reference:** LINE private API protocol derived from [fadhiilrachman/line-py](https://github.com/fadhiilrachman/line-py) — check `linepy/client.py`, `linepy/talk.py`, `linepy/auth.py` for endpoints and headers.

---

## Chunk 1: Project Setup + Transport Layer

### Task 1: Initialize Go Module

**Files:**
- Create: `go.mod`
- Create: `go.sum`

- [ ] **Step 1: Initialize module**

```bash
cd /Users/vysina/my_workspace/linebotgo
go mod init github.com/vysina/linebotgo
```

Expected: `go.mod` created with `module github.com/vysina/linebotgo` and `go 1.21`

- [ ] **Step 2: Add dependencies**

```bash
go get nhooyr.io/websocket@latest
go get github.com/stretchr/testify@latest
go mod tidy
```

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: initialize go module"
```

---

### Task 2: Transport Layer — HTTP Client

**Files:**
- Create: `transport/http.go`
- Create: `transport/http_test.go`

LINE's private API requires specific headers to look like a mobile client. This layer sets them.

- [ ] **Step 1: Write failing test**

```go
// transport/http_test.go
package transport_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/vysina/linebotgo/transport"
)

func TestHTTPClientSetsLINEHeaders(t *testing.T) {
    var gotHeaders http.Header
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        gotHeaders = r.Header
        w.WriteHeader(200)
    }))
    defer srv.Close()

    c := transport.NewHTTPClient()
    req, _ := http.NewRequest("GET", srv.URL, nil)
    c.Do(req)

    assert.Equal(t, "jp.naver.line.android", gotHeaders.Get("X-Line-Application"))
    assert.NotEmpty(t, gotHeaders.Get("User-Agent"))
}
```

- [ ] **Step 2: Run test — verify it fails**

```bash
go test ./transport/ -run TestHTTPClientSetsLINEHeaders -v
```

Expected: FAIL — `transport` package does not exist

- [ ] **Step 3: Implement `transport/http.go`**

```go
// transport/http.go
package transport

import (
    "net/http"
    "time"
)

const (
    // LINE Android client identifiers — update via update-line-api skill if LINE changes these
    LineApplication  = "jp.naver.line.android"
    LineUserAgent    = "Line/13.18.1"
    LineSystemName   = "Android OS"
    LineAppVersion   = "13.18.1"
)

// HTTPClient wraps net/http.Client with LINE-specific defaults.
type HTTPClient struct {
    inner *http.Client
}

// NewHTTPClient returns an HTTPClient configured with LINE mobile client headers.
func NewHTTPClient() *HTTPClient {
    return &HTTPClient{
        inner: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// Do executes the request after injecting LINE client headers.
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
    req.Header.Set("X-Line-Application", LineApplication)
    req.Header.Set("User-Agent", LineUserAgent)
    req.Header.Set("X-Line-System-Name", LineSystemName)
    req.Header.Set("Content-Type", "application/x-thrift")
    req.Header.Set("Accept", "application/x-thrift")
    return c.inner.Do(req)
}
```

- [ ] **Step 4: Run test — verify it passes**

```bash
go test ./transport/ -run TestHTTPClientSetsLINEHeaders -v
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add transport/
git commit -m "feat: add transport layer with LINE HTTP client headers"
```

---

### Task 3: Transport Layer — Thrift Helpers

LINE uses Apache Thrift (binary protocol) for most API calls.

**Files:**
- Create: `transport/thrift.go`
- Create: `transport/thrift_test.go`

- [ ] **Step 1: Add Thrift dependency**

```bash
go get github.com/apache/thrift@latest
go mod tidy
```

- [ ] **Step 2: Write failing test**

```go
// transport/thrift_test.go
package transport_test

import (
    "bytes"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/vysina/linebotgo/transport"
)

func TestThriftReaderWriter(t *testing.T) {
    buf := &bytes.Buffer{}
    w := transport.NewThriftWriter(buf)
    err := w.WriteString("hello")
    assert.NoError(t, err)

    r := transport.NewThriftReader(bytes.NewReader(buf.Bytes()))
    got, err := r.ReadString()
    assert.NoError(t, err)
    assert.Equal(t, "hello", got)
}
```

- [ ] **Step 3: Run — verify it fails**

```bash
go test ./transport/ -run TestThriftReaderWriter -v
```

- [ ] **Step 4: Implement `transport/thrift.go`**

```go
// transport/thrift.go
package transport

import (
    "io"

    "github.com/apache/thrift/lib/go/thrift"
)

// ThriftWriter wraps a Thrift binary protocol writer.
type ThriftWriter struct {
    proto *thrift.TBinaryProtocol
}

func NewThriftWriter(w io.Writer) *ThriftWriter {
    trans := thrift.NewStreamTransportW(w)
    return &ThriftWriter{proto: thrift.NewTBinaryProtocolTransport(trans)}
}

func (w *ThriftWriter) WriteString(s string) error {
    return w.proto.WriteString(s)
}

// ThriftReader wraps a Thrift binary protocol reader.
type ThriftReader struct {
    proto *thrift.TBinaryProtocol
}

func NewThriftReader(r io.Reader) *ThriftReader {
    trans := thrift.NewStreamTransportR(r)
    return &ThriftReader{proto: thrift.NewTBinaryProtocolTransport(trans)}
}

func (r *ThriftReader) ReadString() (string, error) {
    return r.proto.ReadString()
}
```

- [ ] **Step 5: Run — verify it passes**

```bash
go test ./transport/ -v
```

Expected: all PASS

- [ ] **Step 6: Commit**

```bash
git add transport/thrift.go transport/thrift_test.go go.mod go.sum
git commit -m "feat: add thrift read/write helpers to transport layer"
```

---

## Chunk 2: Config + Auth Layer

### Task 4: Config — Token File Persistence

**Files:**
- Create: `config/config.go`
- Create: `config/config_test.go`

- [ ] **Step 1: Write failing test**

```go
// config/config_test.go
package config_test

import (
    "os"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/vysina/linebotgo/config"
)

func TestSaveAndLoadToken(t *testing.T) {
    f, _ := os.CreateTemp("", "linebotgo-test-*.json")
    f.Close()
    defer os.Remove(f.Name())

    cfg := &config.Config{
        AuthToken:    "test-token",
        RefreshToken: "test-refresh",
        ExpiresAt:    time.Now().Add(24 * time.Hour),
        DeviceID:     "device-123",
    }

    err := config.Save(f.Name(), cfg)
    assert.NoError(t, err)

    loaded, err := config.Load(f.Name())
    assert.NoError(t, err)
    assert.Equal(t, "test-token", loaded.AuthToken)
    assert.Equal(t, "device-123", loaded.DeviceID)
}

func TestLoadMissingFile(t *testing.T) {
    _, err := config.Load("/tmp/does-not-exist-linebotgo.json")
    assert.ErrorIs(t, err, config.ErrNoConfig)
}
```

- [ ] **Step 2: Run — verify it fails**

```bash
go test ./config/ -v
```

- [ ] **Step 3: Implement `config/config.go`**

```go
// config/config.go
package config

import (
    "encoding/json"
    "errors"
    "os"
    "time"
)

var ErrNoConfig = errors.New("config file not found")

// Config holds persisted auth state between sessions.
type Config struct {
    AuthToken    string    `json:"auth_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
    DeviceID     string    `json:"device_id"`
    Certificate  string    `json:"certificate"`
}

// Save writes cfg to a JSON file at path.
func Save(path string, cfg *Config) error {
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0600)
}

// Load reads a Config from the JSON file at path.
// Returns ErrNoConfig if the file does not exist.
func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if errors.Is(err, os.ErrNotExist) {
        return nil, ErrNoConfig
    }
    if err != nil {
        return nil, err
    }
    var cfg Config
    return &cfg, json.Unmarshal(data, &cfg)
}
```

- [ ] **Step 4: Run — verify it passes**

```bash
go test ./config/ -v
```

- [ ] **Step 5: Commit**

```bash
git add config/
git commit -m "feat: add config layer for token persistence"
```

---

### Task 5: Auth — Error Types + Token

**Files:**
- Create: `auth/errors.go`
- Create: `auth/token.go`
- Create: `auth/token_test.go`

- [ ] **Step 1: Write failing test**

```go
// auth/token_test.go
package auth_test

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/vysina/linebotgo/auth"
)

func TestTokenExpiry(t *testing.T) {
    expired := &auth.Token{
        AuthToken: "tok",
        ExpiresAt: time.Now().Add(-1 * time.Hour),
    }
    assert.True(t, expired.IsExpired())

    valid := &auth.Token{
        AuthToken: "tok",
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
    assert.False(t, valid.IsExpired())
}
```

- [ ] **Step 2: Run — verify it fails**

```bash
go test ./auth/ -run TestTokenExpiry -v
```

- [ ] **Step 3: Implement auth/errors.go and auth/token.go**

```go
// auth/errors.go
package auth

import "errors"

var (
    ErrTokenExpired  = errors.New("auth token expired")
    ErrInvalidToken  = errors.New("auth token is invalid")
    ErrRateLimited   = errors.New("rate limited by LINE")
    ErrLoginBlocked  = errors.New("login blocked by LINE")
    ErrLoginFailed   = errors.New("login failed: invalid credentials")
)
```

```go
// auth/token.go
package auth

import "time"

// Token holds an active LINE session.
type Token struct {
    AuthToken    string
    RefreshToken string
    ExpiresAt    time.Time
    DeviceID     string
    Certificate  string
}

// IsExpired returns true if the token is past its expiry time.
func (t *Token) IsExpired() bool {
    return time.Now().After(t.ExpiresAt)
}
```

- [ ] **Step 4: Run — verify it passes**

```bash
go test ./auth/ -run TestTokenExpiry -v
```

- [ ] **Step 5: Commit**

```bash
git add auth/
git commit -m "feat: add auth error types and token model"
```

---

### Task 6: Auth — Email/Password Login

**Files:**
- Create: `auth/email.go`
- Create: `auth/email_test.go`

Reference: `linepy/auth.py` — `loginWithIdentityCredential` POST to `https://gw.line.naver.jp/api/v4/TalkService/loginWithIdentityCredential`

- [ ] **Step 1: Write failing test with mock HTTP server**

```go
// auth/email_test.go
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
        json.NewEncoder(w).Encode(map[string]interface{}{
            "authToken":    "mock-token",
            "refreshToken": "mock-refresh",
            "certificate":  "mock-cert",
        })
    }))
    defer srv.Close()

    client := transport.NewHTTPClient()
    loginURL := srv.URL + "/api/v4/TalkService/loginWithIdentityCredential"
    tok, err := auth.EmailLogin(client, loginURL, "test@example.com", "password123")

    assert.NoError(t, err)
    assert.Equal(t, "mock-token", tok.AuthToken)
    assert.Equal(t, "mock-refresh", tok.RefreshToken)
}

func TestEmailLogin_BadCredentials(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusUnauthorized)
    }))
    defer srv.Close()

    client := transport.NewHTTPClient()
    _, err := auth.EmailLogin(client, srv.URL+"/login", "bad@example.com", "wrong")
    assert.ErrorIs(t, err, auth.ErrLoginFailed)
}
```

- [ ] **Step 2: Run — verify it fails**

```bash
go test ./auth/ -run TestEmailLogin -v
```

- [ ] **Step 3: Implement `auth/email.go`**

```go
// auth/email.go
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
// loginURL should be the full endpoint URL (allows override in tests).
func EmailLogin(c *transport.HTTPClient, loginURL, email, password string) (*Token, error) {
    body := fmt.Sprintf(`{"identifier":{"type":"email","value":"%s"},"secret":"%s"}`, email, password)
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

    if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
        return nil, ErrLoginFailed
    }
    if resp.StatusCode == http.StatusTooManyRequests {
        return nil, ErrRateLimited
    }
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
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
        ExpiresAt:    time.Now().Add(24 * time.Hour), // LINE tokens expire ~24h
    }, nil
}
```

- [ ] **Step 4: Run — verify tests pass**

```bash
go test ./auth/ -run TestEmailLogin -v
```

- [ ] **Step 5: Commit**

```bash
git add auth/email.go auth/email_test.go
git commit -m "feat: add email/password login flow"
```

---

### Task 7: Auth — QR Code Login

**Files:**
- Create: `auth/qrcode.go`
- Create: `auth/qrcode_test.go`

QR login: fetch QR session → display QR in terminal → poll for scan confirmation.

- [ ] **Step 1: Write failing test**

```go
// auth/qrcode_test.go
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

func TestQRLogin_IssuesQR(t *testing.T) {
    callCount := 0
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        callCount++
        switch callCount {
        case 1: // create session
            json.NewEncoder(w).Encode(map[string]string{"session": "qr-session-id", "url": "https://line.me/R/qr?code=abc"})
        case 2: // poll — first: not scanned yet
            w.WriteHeader(http.StatusAccepted)
        case 3: // poll — scanned
            json.NewEncoder(w).Encode(map[string]string{"authToken": "qr-token", "refreshToken": "qr-refresh"})
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
}
```

- [ ] **Step 2: Run — verify it fails**

```bash
go test ./auth/ -run TestQRLogin -v
```

- [ ] **Step 3: Implement `auth/qrcode.go`**

```go
// auth/qrcode.go
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

// NewQRLogin creates a QRLogin instance. URLs allow override in tests.
func NewQRLogin(c *transport.HTTPClient, createURL, pollURL string) *QRLogin {
    return &QRLogin{client: c, createURL: createURL, pollURL: pollURL}
}

// Login fetches a QR code, calls displayFn with the URL, then polls until scanned.
func (q *QRLogin) Login(displayFn func(qrURL string)) (*Token, error) {
    // Step 1: create session
    req, _ := http.NewRequest("POST", q.createURL, nil)
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

    // Step 2: poll until scanned (max 2 minutes)
    deadline := time.Now().Add(2 * time.Minute)
    for time.Now().Before(deadline) {
        req, _ := http.NewRequest("GET", q.pollURL+"?session="+session.Session, nil)
        resp, err := q.client.Do(req)
        if err != nil {
            return nil, err
        }

        if resp.StatusCode == http.StatusAccepted {
            resp.Body.Close()
            time.Sleep(2 * time.Second)
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
    return nil, fmt.Errorf("QR login timed out")
}
```

- [ ] **Step 4: Run — verify tests pass**

```bash
go test ./auth/ -v
```

- [ ] **Step 5: Commit**

```bash
git add auth/qrcode.go auth/qrcode_test.go
git commit -m "feat: add QR code login flow"
```

---

## Chunk 3: API Definitions + Core Client

### Task 8: API Layer — Endpoints and Types

**Files:**
- Create: `api/endpoints.go` ← **main target for update-line-api skill**
- Create: `api/types.go`

- [ ] **Step 1: Create `api/endpoints.go`**

```go
// api/endpoints.go
// IMPORTANT: This file is the primary target for the update-line-api Claude skill.
// When LINE changes API endpoints or parameters, update this file.
package api

const (
    BaseURL = "https://gw.line.naver.jp"

    // Auth
    EndpointEmailLogin   = BaseURL + "/api/v4/TalkService/loginWithIdentityCredential"
    EndpointQRCreate     = BaseURL + "/api/v4/TalkService/getAuthQrCode"
    EndpointQRPoll       = BaseURL + "/api/v4/TalkService/checkAuthQrCodeVerified"
    EndpointRefreshToken = BaseURL + "/api/v4/TalkService/refreshToken"

    // Talk / Messaging
    EndpointFetchOps    = BaseURL + "/SYNC/v3/TalkService/fetchOps"
    EndpointSendMessage = BaseURL + "/S4/TalkService/sendMessage"
    EndpointGetContacts = BaseURL + "/S4/TalkService/getAllContactIds"
    EndpointGetProfile  = BaseURL + "/S4/TalkService/getProfile"

    // OpenChat
    EndpointOpenChatSearch    = BaseURL + "/api/v1/square/search"
    EndpointOpenChatJoin      = BaseURL + "/api/v1/square/join"
    EndpointOpenChatLeave     = BaseURL + "/api/v1/square/leave"
    EndpointOpenChatSend      = BaseURL + "/api/v1/square/message/send"
    EndpointOpenChatMembers   = BaseURL + "/api/v1/square/members"
    EndpointOpenChatKick      = BaseURL + "/api/v1/square/members/kick"
    EndpointOpenChatBan       = BaseURL + "/api/v1/square/members/ban"
    EndpointOpenChatPin       = BaseURL + "/api/v1/square/message/pin"
    EndpointOpenChatInfo      = BaseURL + "/api/v1/square/info"
)
```

- [ ] **Step 2: Create `api/types.go`**

```go
// api/types.go
package api

import "time"

// Message represents a LINE message (talk or OpenChat).
type Message struct {
    ID          string
    Type        MessageType
    From        string // sender MID
    To          string // recipient/group ID
    Text        string
    ContentURL  string
    CreatedAt   time.Time
}

type MessageType string

const (
    MessageTypeText     MessageType = "NONE" // LINE uses NONE for text
    MessageTypeSticker  MessageType = "STICKER"
    MessageTypeImage    MessageType = "IMAGE"
    MessageTypeVideo    MessageType = "VIDEO"
    MessageTypeAudio    MessageType = "AUDIO"
    MessageTypeLocation MessageType = "LOCATION"
    MessageTypeContact  MessageType = "CONTACT"
)

// Contact represents a LINE contact/friend.
type Contact struct {
    MID         string
    DisplayName string
    StatusMessage string
    PictureURL  string
}

// OpenChatInfo represents an OpenChat room.
type OpenChatInfo struct {
    ID          string
    Name        string
    Description string
    MemberCount int
}

// OpenChatMember represents a member of an OpenChat room.
type OpenChatMember struct {
    MID         string
    DisplayName string
    Role        string // "ADMIN" | "MEMBER"
}
```

- [ ] **Step 3: Verify it compiles**

```bash
go build ./api/
```

- [ ] **Step 4: Commit**

```bash
git add api/
git commit -m "feat: add API endpoint definitions and types"
```

---

### Task 9: Core Client

**Files:**
- Create: `client/client.go`
- Create: `client/client_test.go`
- Create: `client/message.go`
- Create: `client/message_test.go`
- Create: `client/contact.go`
- Create: `client/poll.go`

- [ ] **Step 1: Write failing test for client creation**

```go
// client/client_test.go
package client_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/vysina/linebotgo/auth"
    "github.com/vysina/linebotgo/client"
    "time"
)

func TestNewClientWithToken(t *testing.T) {
    tok := &auth.Token{
        AuthToken: "test-token",
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
    c := client.New(client.WithToken(tok))
    assert.NotNil(t, c)
    assert.Equal(t, "test-token", c.AuthToken())
}
```

- [ ] **Step 2: Run — verify it fails**

```bash
go test ./client/ -run TestNewClientWithToken -v
```

- [ ] **Step 3: Implement `client/client.go`**

```go
// client/client.go
package client

import (
    "github.com/vysina/linebotgo/auth"
    "github.com/vysina/linebotgo/transport"
)

// Client is the main LINE API client.
type Client struct {
    token     *auth.Token
    http      *transport.HTTPClient
    msgHandlers     []func(*Message)
    ocMsgHandlers   map[string][]func(*OpenChatMessage)
}

// Option configures a Client.
type Option func(*Client)

// WithToken creates a client from an existing auth token.
func WithToken(t *auth.Token) Option {
    return func(c *Client) { c.token = t }
}

// New creates a new Client with the given options.
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

// Message is a received LINE message (re-exported from api for convenience).
type Message struct {
    ID    string
    From  string
    To    string
    Text  string
    Type  string
}

// OpenChatMessage is a received OpenChat message.
type OpenChatMessage struct {
    ChatID string
    From   string
    Text   string
}
```

- [ ] **Step 4: Write failing test for SendText**

```go
// client/message_test.go
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

func TestSendText(t *testing.T) {
    var gotBody map[string]interface{}
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewDecoder(r.Body).Decode(&gotBody)
        w.WriteHeader(200)
    }))
    defer srv.Close()

    tok := &auth.Token{AuthToken: "test-token", ExpiresAt: time.Now().Add(24 * time.Hour)}
    c := client.New(client.WithToken(tok), client.WithSendURL(srv.URL+"/send"))

    err := c.SendText("friend-mid", "Hello!")
    assert.NoError(t, err)
    assert.Equal(t, "friend-mid", gotBody["to"])
    assert.Equal(t, "Hello!", gotBody["text"])
}
```

- [ ] **Step 5: Run — verify it fails**

```bash
go test ./client/ -run TestSendText -v
```

- [ ] **Step 6: Implement `client/message.go` and `WithSendURL` option**

```go
// client/message.go
package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/vysina/linebotgo/api"
)

// sendURL allows override in tests; defaults to api.EndpointSendMessage.
func (c *Client) resolveSendURL() string {
    if c.sendURL != "" {
        return c.sendURL
    }
    return api.EndpointSendMessage
}

// SendText sends a text message to a LINE user or group.
func (c *Client) SendText(to, text string) error {
    payload := map[string]interface{}{
        "to":   to,
        "text": text,
        "type": "NONE",
    }
    return c.postJSON(c.resolveSendURL(), payload)
}

// OnMessage registers a handler called for each incoming message.
func (c *Client) OnMessage(fn func(*Message)) {
    c.msgHandlers = append(c.msgHandlers, fn)
}

func (c *Client) postJSON(url string, payload interface{}) error {
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
        return fmt.Errorf("LINE API error: %d", resp.StatusCode)
    }
    return nil
}
```

Add `sendURL` field and `WithSendURL` option to `client/client.go`:

```go
// in Client struct:
sendURL string

// new option:
func WithSendURL(url string) Option {
    return func(c *Client) { c.sendURL = url }
}
```

- [ ] **Step 7: Implement `client/contact.go`**

```go
// client/contact.go
package client

import (
    "encoding/json"
    "net/http"

    "github.com/vysina/linebotgo/api"
)

// GetProfile returns profile information for a given MID.
func (c *Client) GetProfile(mid string) (*api.Contact, error) {
    req, _ := http.NewRequest("GET", api.EndpointGetProfile+"?mid="+mid, nil)
    req.Header.Set("X-Line-Access", c.AuthToken())
    resp, err := c.http.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    var contact api.Contact
    return &contact, json.NewDecoder(resp.Body).Decode(&contact)
}
```

- [ ] **Step 8: Implement `client/poll.go`** (long polling for incoming messages)

```go
// client/poll.go
package client

import (
    "context"
    "encoding/json"
    "net/http"
    "time"

    "github.com/vysina/linebotgo/api"
)

// Listen starts long polling and calls registered OnMessage handlers.
// It blocks until ctx is cancelled.
func (c *Client) Listen(ctx context.Context) error {
    revision := int64(0)
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        ops, err := c.fetchOps(ctx, revision)
        if err != nil {
            time.Sleep(5 * time.Second)
            continue
        }

        for _, op := range ops {
            if op.Message != nil {
                for _, h := range c.msgHandlers {
                    h(op.Message)
                }
            }
            if op.Revision > revision {
                revision = op.Revision
            }
        }
    }
}

type operation struct {
    Revision int64
    Message  *Message
}

func (c *Client) fetchOps(ctx context.Context, revision int64) ([]operation, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET",
        api.EndpointFetchOps, nil)
    req.Header.Set("X-Line-Access", c.AuthToken())
    q := req.URL.Query()
    q.Set("revision", fmt.Sprintf("%d", revision))
    q.Set("count", "50")
    req.URL.RawQuery = q.Encode()

    resp, err := c.http.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result struct {
        Operations []struct {
            Revision int64 `json:"revision"`
            Message  *Message `json:"message,omitempty"`
        } `json:"operations"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    ops := make([]operation, len(result.Operations))
    for i, o := range result.Operations {
        ops[i] = operation{Revision: o.Revision, Message: o.Message}
    }
    return ops, nil
}
```

Fix the missing `fmt` import — add `"fmt"` to `client/poll.go` imports.

- [ ] **Step 9: Run all client tests**

```bash
go test ./client/ -v
```

Expected: all PASS

- [ ] **Step 10: Commit**

```bash
git add client/
git commit -m "feat: add core client with send, contacts, long polling"
```

---

## Chunk 4: OpenChat Layer

### Task 10: OpenChat Core Operations

**Files:**
- Create: `openchat/openchat.go`
- Create: `openchat/openchat_test.go`
- Create: `openchat/member.go`
- Create: `openchat/admin.go`
- Create: `openchat/search.go`

- [ ] **Step 1: Write failing test for search and join**

```go
// openchat/openchat_test.go
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
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
    }))
    defer srv.Close()

    oc := openchat.New(newTestClient(t, srv.URL))
    err := oc.Join("chat-1")
    assert.NoError(t, err)
}
```

- [ ] **Step 2: Run — verify it fails**

```bash
go test ./openchat/ -v
```

- [ ] **Step 3: Add `WithBaseURL` option to client**

In `client/client.go`, add `baseURL string` field and option:
```go
func WithBaseURL(url string) Option {
    return func(c *Client) { c.baseURL = url }
}
```

Update `resolveSendURL()` and all URL-resolving methods to prepend `c.baseURL` if set.

- [ ] **Step 4: Implement `openchat/openchat.go`**

```go
// openchat/openchat.go
package openchat

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/vysina/linebotgo/api"
    "github.com/vysina/linebotgo/client"
)

// OpenChat provides LINE OpenChat operations.
type OpenChat struct {
    c *client.Client
}

// New creates an OpenChat accessor from an existing client.
func New(c *client.Client) *OpenChat {
    return &OpenChat{c: c}
}

// Search finds OpenChat rooms matching the query.
func (oc *OpenChat) Search(query string) ([]api.OpenChatInfo, error) {
    req, _ := http.NewRequest("GET", oc.url(api.EndpointOpenChatSearch), nil)
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
    return result.Chats, json.NewDecoder(resp.Body).Decode(&result)
}

// Join joins an OpenChat room by ID.
func (oc *OpenChat) Join(chatID string) error {
    return oc.c.PostJSON(oc.url(api.EndpointOpenChatJoin), map[string]string{"chatId": chatID})
}

// Leave leaves an OpenChat room.
func (oc *OpenChat) Leave(chatID string) error {
    return oc.c.PostJSON(oc.url(api.EndpointOpenChatLeave), map[string]string{"chatId": chatID})
}

// SendText sends a text message to an OpenChat room.
func (oc *OpenChat) SendText(chatID, text string) error {
    return oc.c.PostJSON(oc.url(api.EndpointOpenChatSend), map[string]string{
        "chatId": chatID,
        "text":   text,
    })
}

func (oc *OpenChat) url(endpoint string) string {
    if base := oc.c.BaseURL(); base != "" {
        return base + "/" + endpoint[len("https://gw.line.naver.jp"):]
    }
    return endpoint
}

func (oc *OpenChat) OnMessage(chatID string, fn func(*client.OpenChatMessage)) {
    oc.c.RegisterOpenChatHandler(chatID, fn)
}
```

Expose helpers from `client.Client`: `HTTPClient()`, `PostJSON()` (make exported), `BaseURL()`, `RegisterOpenChatHandler()`.

- [ ] **Step 5: Implement `openchat/member.go`**

```go
// openchat/member.go
package openchat

import (
    "encoding/json"
    "net/http"

    "github.com/vysina/linebotgo/api"
)

// GetMembers returns members of an OpenChat room.
func (oc *OpenChat) GetMembers(chatID string) ([]api.OpenChatMember, error) {
    req, _ := http.NewRequest("GET", oc.url(api.EndpointOpenChatMembers)+"?chatId="+chatID, nil)
    req.Header.Set("X-Line-Access", oc.c.AuthToken())
    resp, err := oc.c.HTTPClient().Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    var result struct {
        Members []api.OpenChatMember `json:"members"`
    }
    return result.Members, json.NewDecoder(resp.Body).Decode(&result)
}

// Kick removes a member from an OpenChat room.
func (oc *OpenChat) Kick(chatID, memberMID string) error {
    return oc.c.PostJSON(oc.url(api.EndpointOpenChatKick), map[string]string{
        "chatId": chatID, "mid": memberMID,
    })
}

// Ban bans a member from an OpenChat room.
func (oc *OpenChat) Ban(chatID, memberMID string) error {
    return oc.c.PostJSON(oc.url(api.EndpointOpenChatBan), map[string]string{
        "chatId": chatID, "mid": memberMID,
    })
}
```

- [ ] **Step 6: Implement `openchat/admin.go`**

```go
// openchat/admin.go
package openchat

import (
    "encoding/json"
    "net/http"

    "github.com/vysina/linebotgo/api"
)

// PinMessage pins a message in an OpenChat room.
func (oc *OpenChat) PinMessage(chatID, messageID string) error {
    return oc.c.PostJSON(oc.url(api.EndpointOpenChatPin), map[string]string{
        "chatId": chatID, "messageId": messageID,
    })
}

// SetNotification enables or disables notifications for an OpenChat room.
func (oc *OpenChat) SetNotification(chatID string, enabled bool) error {
    return oc.c.PostJSON(oc.url(api.EndpointOpenChatInfo)+"/notification", map[string]interface{}{
        "chatId": chatID, "enabled": enabled,
    })
}

// GetChatInfo returns info about an OpenChat room.
func (oc *OpenChat) GetChatInfo(chatID string) (*api.OpenChatInfo, error) {
    req, _ := http.NewRequest("GET", oc.url(api.EndpointOpenChatInfo)+"?chatId="+chatID, nil)
    req.Header.Set("X-Line-Access", oc.c.AuthToken())
    resp, err := oc.c.HTTPClient().Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    var info api.OpenChatInfo
    return &info, json.NewDecoder(resp.Body).Decode(&info)
}
```

- [ ] **Step 7: Run all tests**

```bash
go test ./... -v
```

Expected: all PASS

- [ ] **Step 8: Commit**

```bash
git add openchat/
git commit -m "feat: add OpenChat layer (search, join, messaging, members, admin)"
```

---

## Chunk 5: HTTP/WebSocket Server

### Task 11: Server Setup + REST Endpoints

**Files:**
- Create: `server/server.go`
- Create: `server/rest.go`
- Create: `server/rest_test.go`
- Create: `server/middleware/logging.go`

- [ ] **Step 1: Add WebSocket dependency**

```bash
go get nhooyr.io/websocket@latest
go mod tidy
```

- [ ] **Step 2: Write failing test for REST send endpoint**

```go
// server/rest_test.go
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

type mockClient struct {
    lastTo   string
    lastText string
}

func (m *mockClient) SendText(to, text string) error {
    m.lastTo = to
    m.lastText = text
    return nil
}

func TestRESTSendMessage(t *testing.T) {
    mc := &mockClient{}
    srv := server.New(mc, server.Config{Port: 0})
    handler := srv.Handler()

    body, _ := json.Marshal(map[string]string{"to": "user-mid", "text": "hi"})
    req := httptest.NewRequest("POST", "/api/message/send", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    handler.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "user-mid", mc.lastTo)
    assert.Equal(t, "hi", mc.lastText)
}
```

- [ ] **Step 3: Run — verify it fails**

```bash
go test ./server/ -run TestRESTSendMessage -v
```

- [ ] **Step 4: Implement `server/server.go`**

```go
// server/server.go
package server

import (
    "fmt"
    "net/http"
)

// Sender is the interface the server uses to send messages.
// Keeps server decoupled from client implementation.
type Sender interface {
    SendText(to, text string) error
}

// Config holds server configuration.
type Config struct {
    Port int
}

// Server wraps a Sender and serves REST + WebSocket.
type Server struct {
    sender  Sender
    cfg     Config
    hub     *wsHub
    mux     *http.ServeMux
}

// New creates a new Server.
func New(sender Sender, cfg Config) *Server {
    s := &Server{
        sender: sender,
        cfg:    cfg,
        hub:    newWSHub(),
        mux:    http.NewServeMux(),
    }
    s.registerRoutes()
    return s
}

// Handler returns the HTTP handler (for testing without starting a full server).
func (s *Server) Handler() http.Handler {
    return s.mux
}

// Start starts the HTTP server.
func (s *Server) Start() error {
    addr := fmt.Sprintf(":%d", s.cfg.Port)
    return http.ListenAndServe(addr, s.mux)
}

func (s *Server) registerRoutes() {
    s.mux.HandleFunc("POST /api/message/send", s.handleSendMessage)
    s.mux.HandleFunc("GET /api/contacts", s.handleGetContacts)
    s.mux.HandleFunc("GET /ws/events", s.handleWebSocket)
}
```

- [ ] **Step 5: Implement `server/rest.go`**

```go
// server/rest.go
package server

import (
    "encoding/json"
    "net/http"
)

func (s *Server) handleSendMessage(w http.ResponseWriter, r *http.Request) {
    var req struct {
        To   string `json:"to"`
        Text string `json:"text"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }
    if err := s.sender.SendText(req.To, req.Text); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetContacts(w http.ResponseWriter, r *http.Request) {
    // Placeholder — extend with ContactGetter interface when needed
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"contacts": []interface{}{}})
}
```

- [ ] **Step 6: Run REST tests**

```bash
go test ./server/ -run TestRESTSendMessage -v
```

Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add server/server.go server/rest.go server/rest_test.go
git commit -m "feat: add HTTP server with REST send-message endpoint"
```

---

### Task 12: WebSocket Event Broadcasting

**Files:**
- Create: `server/ws.go`
- Create: `server/ws_test.go`

- [ ] **Step 1: Write failing test for WebSocket**

```go
// server/ws_test.go
package server_test

import (
    "context"
    "encoding/json"
    "testing"

    "github.com/stretchr/testify/assert"
    "nhooyr.io/websocket"
    "nhooyr.io/websocket/wsjson"
    "net/http/httptest"
    "github.com/vysina/linebotgo/server"
)

func TestWebSocketReceivesEvent(t *testing.T) {
    mc := &mockClient{}
    srv := server.New(mc, server.Config{})
    ts := httptest.NewServer(srv.Handler())
    defer ts.Close()

    wsURL := "ws" + ts.URL[4:] + "/ws/events"
    ctx := context.Background()
    conn, _, err := websocket.Dial(ctx, wsURL, nil)
    assert.NoError(t, err)
    defer conn.Close(websocket.StatusNormalClosure, "")

    // Broadcast an event from server side
    srv.Broadcast(server.Event{Type: "message", Data: map[string]string{"text": "hello"}})

    var received server.Event
    err = wsjson.Read(ctx, conn, &received)
    assert.NoError(t, err)
    assert.Equal(t, "message", received.Type)
}
```

- [ ] **Step 2: Run — verify it fails**

```bash
go test ./server/ -run TestWebSocketReceivesEvent -v
```

- [ ] **Step 3: Implement `server/ws.go`**

```go
// server/ws.go
package server

import (
    "context"
    "encoding/json"
    "net/http"
    "sync"

    "nhooyr.io/websocket"
    "nhooyr.io/websocket/wsjson"
)

// Event is a real-time event broadcast over WebSocket.
type Event struct {
    Type string      `json:"type"`
    Data interface{} `json:"data"`
}

type wsHub struct {
    mu    sync.Mutex
    conns map[*websocket.Conn]struct{}
}

func newWSHub() *wsHub {
    return &wsHub{conns: make(map[*websocket.Conn]struct{})}
}

func (h *wsHub) add(c *websocket.Conn) {
    h.mu.Lock()
    h.conns[c] = struct{}{}
    h.mu.Unlock()
}

func (h *wsHub) remove(c *websocket.Conn) {
    h.mu.Lock()
    delete(h.conns, c)
    h.mu.Unlock()
}

func (h *wsHub) broadcast(e Event) {
    h.mu.Lock()
    defer h.mu.Unlock()
    for c := range h.conns {
        wsjson.Write(context.Background(), c, e) // best-effort
    }
}

// Broadcast sends an event to all connected WebSocket clients.
func (s *Server) Broadcast(e Event) {
    s.hub.broadcast(e)
}

// BroadcastMessage is a convenience method called when the client receives a message.
func (s *Server) BroadcastMessage(from, to, text string) {
    s.Broadcast(Event{
        Type: "message",
        Data: map[string]string{"from": from, "to": to, "text": text},
    })
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
        InsecureSkipVerify: true, // for local/dev use
    })
    if err != nil {
        return
    }
    s.hub.add(conn)
    defer s.hub.remove(conn)

    // Keep connection alive until client disconnects
    ctx := conn.CloseRead(context.Background())
    <-ctx.Done()
}
```

- [ ] **Step 4: Run all server tests**

```bash
go test ./server/ -v
```

Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add server/ws.go server/ws_test.go
git commit -m "feat: add WebSocket broadcast hub"
```

---

## Chunk 6: Examples + Claude Skill

### Task 13: Simple Echo Bot Example

**Files:**
- Create: `examples/simple-bot/main.go`

- [ ] **Step 1: Create `examples/simple-bot/main.go`**

```go
// examples/simple-bot/main.go
// Simple echo bot — echoes every received message back to sender.
package main

import (
    "context"
    "log"
    "os"

    "github.com/vysina/linebotgo/api"
    "github.com/vysina/linebotgo/auth"
    "github.com/vysina/linebotgo/client"
    "github.com/vysina/linebotgo/config"
)

func main() {
    const configPath = "~/.config/linebotgo/config.json"

    var tok *auth.Token

    cfg, err := config.Load(configPath)
    if err == nil {
        tok = &auth.Token{
            AuthToken:    cfg.AuthToken,
            RefreshToken: cfg.RefreshToken,
            ExpiresAt:    cfg.ExpiresAt,
        }
        log.Println("Loaded token from config")
    } else {
        // fallback to email login
        email := os.Getenv("LINE_EMAIL")
        password := os.Getenv("LINE_PASSWORD")
        if email == "" || password == "" {
            log.Fatal("Set LINE_EMAIL and LINE_PASSWORD env vars")
        }
        c := transport.NewHTTPClient()
        tok, err = auth.EmailLogin(c, api.EndpointEmailLogin, email, password)
        if err != nil {
            log.Fatalf("login failed: %v", err)
        }
        config.Save(configPath, &config.Config{
            AuthToken:    tok.AuthToken,
            RefreshToken: tok.RefreshToken,
            ExpiresAt:    tok.ExpiresAt,
        })
        log.Println("Logged in and saved token")
    }

    bot := client.New(client.WithToken(tok))
    bot.OnMessage(func(msg *client.Message) {
        log.Printf("[%s]: %s", msg.From, msg.Text)
        if err := bot.SendText(msg.From, msg.Text); err != nil {
            log.Printf("send error: %v", err)
        }
    })

    log.Println("Bot listening...")
    if err := bot.Listen(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

- [ ] **Step 2: Verify it builds**

```bash
go build ./examples/simple-bot/
```

- [ ] **Step 3: Commit**

```bash
git add examples/simple-bot/
git commit -m "feat: add simple echo bot example"
```

---

### Task 14: API Server Example

**Files:**
- Create: `examples/api-server/main.go`

- [ ] **Step 1: Create `examples/api-server/main.go`**

```go
// examples/api-server/main.go
// Runs linebotgo as an HTTP/WebSocket API server.
package main

import (
    "context"
    "log"
    "os"

    "github.com/vysina/linebotgo/api"
    "github.com/vysina/linebotgo/auth"
    "github.com/vysina/linebotgo/client"
    "github.com/vysina/linebotgo/config"
    "github.com/vysina/linebotgo/server"
    "github.com/vysina/linebotgo/transport"
)

func main() {
    const configPath = "~/.config/linebotgo/config.json"

    cfg, err := config.Load(configPath)
    var tok *auth.Token
    if err == nil {
        tok = &auth.Token{AuthToken: cfg.AuthToken, ExpiresAt: cfg.ExpiresAt}
    } else {
        email := os.Getenv("LINE_EMAIL")
        password := os.Getenv("LINE_PASSWORD")
        c := transport.NewHTTPClient()
        tok, _ = auth.EmailLogin(c, api.EndpointEmailLogin, email, password)
    }

    bot := client.New(client.WithToken(tok))

    srv := server.New(bot, server.Config{Port: 8080})

    // Forward incoming LINE messages to WebSocket clients
    bot.OnMessage(func(msg *client.Message) {
        srv.BroadcastMessage(msg.From, msg.To, msg.Text)
    })

    go func() {
        log.Println("API server listening on :8080")
        log.Fatal(srv.Start())
    }()

    log.Println("Long polling LINE...")
    bot.Listen(context.Background())
}
```

- [ ] **Step 2: Build**

```bash
go build ./examples/api-server/
```

- [ ] **Step 3: Commit**

```bash
git add examples/api-server/
git commit -m "feat: add API server example"
```

---

### Task 15: Claude Skill — update-line-api

**Files:**
- Create: `skills/update-line-api.md`

- [ ] **Step 1: Create `skills/update-line-api.md`**

````markdown
---
name: update-line-api
description: Update linebotgo's private LINE API definitions when LINE changes their API. Use when encountering API errors, or to analyze network traffic dumps. Updates api/endpoints.go, api/types.go, and related files.
type: project
---

# update-line-api skill

Use this skill when:
- LINE returns unexpected 4xx/5xx errors
- User pastes an error log from linebotgo
- User pastes HTTP traffic (mitmproxy/Charles output)
- User pastes a Thrift/protobuf schema diff
- Running `/update-line-api`

## Workflow

### Step 1 — Read Current State
Read `api/endpoints.go` and `api/types.go` to understand what's currently defined.

### Step 2 — Analyze Input

**If given an error log:**
- Find the endpoint that returned the error
- Identify if it's a URL change, header change, or payload change
- Check `transport/http.go` for header constants

**If given HTTP traffic dump (mitmproxy/Charles):**
- Extract: URL, method, headers, request body, response body
- Compare with current `api/endpoints.go`
- Note any new fields in request/response structs

**If given a Thrift schema diff:**
- Map changed field IDs/types to `api/types.go`
- Identify struct fields to add/remove/rename

### Step 3 — Make Changes

Target files (in priority order):
1. `api/endpoints.go` — URL/path changes
2. `api/types.go` — struct field changes
3. `transport/http.go` — header/version changes (LineApplication, LineUserAgent constants)
4. `auth/email.go` or `auth/qrcode.go` — auth flow changes

### Step 4 — Run Tests
```bash
go test ./... -v
```
Fix any failures before proceeding.

### Step 5 — Update Docs
- Add entry to `CHANGELOG.md` under `## Unreleased`
- If endpoint URLs changed, update README API reference section
- Suggest version bump (patch for fixes, minor for new endpoints)

## Reference
- Original Python impl: https://github.com/fadhiilrachman/line-py
- LINE Thrift IDL can sometimes be found in APK decompiles
- Headers in `transport/http.go` mirror LINE Android app version
````

- [ ] **Step 2: Commit**

```bash
git add skills/
git commit -m "feat: add update-line-api Claude skill"
```

---

### Task 16: Final Verification

- [ ] **Step 1: Run full test suite**

```bash
go test ./... -v
```

Expected: all PASS, no compilation errors

- [ ] **Step 2: Verify all examples build**

```bash
go build ./examples/...
```

- [ ] **Step 3: Run go vet**

```bash
go vet ./...
```

Expected: no issues

- [ ] **Step 4: Final commit**

```bash
git add .
git commit -m "chore: final verification — all tests pass, examples build"
```

---

## Summary

| Chunk | What it builds | Tests |
|-------|---------------|-------|
| 1 | Transport (HTTP client + Thrift) | Unit with mock HTTP server |
| 2 | Config + Auth (email, QR, token) | Unit with mock HTTP server |
| 3 | API definitions + Core client | Unit with mock HTTP server |
| 4 | OpenChat (search, messaging, admin) | Unit with mock HTTP server |
| 5 | HTTP/WebSocket server | Unit + integration with httptest |
| 6 | Examples + Claude skill | Build verification |
