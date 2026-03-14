# linebotgo Design Spec

**Date:** 2026-03-14
**Status:** Approved
**Goal:** Go port of line-py private LINE API library with HTTP/WebSocket server and Claude maintenance skill

---

## Overview

`linebotgo` is a Go library + API server that wraps LINE's private API (not the official LINE Messaging API). It enables programmatic interaction with LINE without the LINE desktop/mobile app. Developers can use it as a Go SDK or run it as an HTTP/WebSocket server to integrate from any language.

---

## Architecture: Layered SDK

```
linebotgo/
├── transport/           ← HTTP client, protobuf encode/decode, TLS
├── auth/                ← Login strategies, token storage, failure handling
├── api/                 ← LINE private API definitions (endpoints, structs)
├── client/              ← High-level Go SDK (entry point for developers)
├── openchat/            ← OpenChat (formerly LINE Square) features
├── config/              ← Config/token file persistence
├── server/              ← HTTP REST + WebSocket server
│   ├── server.go
│   ├── rest.go
│   ├── ws.go
│   └── middleware/
├── examples/
│   ├── simple-bot/      ← Basic echo bot (Go SDK usage)
│   ├── openchat-bot/    ← OpenChat bot template
│   └── api-server/      ← Run as HTTP/WebSocket API server
└── skills/
    └── update-line-api.md  ← Claude skill for API maintenance
```

### Data Flow
```
LINE servers ──long poll──▶ client/ ──▶ server/ws.go ──▶ WebSocket consumers
                                    └──▶ server/rest.go ◀── HTTP requests
Developer Go code ──────────────▶ client/ (direct SDK usage)

All layers: client/ → api/ (definitions) → transport/ → LINE servers
                              ↑
                         auth/ (token injected)
```

---

## Authentication & Token Management

### Login Methods
- **Email/Password** — POST to LINE's loginWithIdentityCredential endpoint
- **QR Code** — fetch QR, display in terminal, poll for scan confirmation

### Token Lifecycle
```
Login → auth token + refresh token
      → save to config file
      → next run: load token, skip login

Token expiry → auto refresh
Refresh fail → re-login
Rate limited → exponential backoff
Login blocked → surface clear error to user
```

### Config File (`~/.config/linebotgo/config.json`)
```json
{
  "auth_token": "...",
  "refresh_token": "...",
  "expires_at": "2026-04-01T00:00:00Z",
  "device_id": "...",
  "certificate": "..."
}
```

### Go SDK Auth
```go
bot, err := linebotgo.New(linebotgo.Config{
    Email:     "xxx@gmail.com",
    Password:  "xxx",
    // or load persisted token:
    TokenFile: "~/.config/linebotgo/config.json",
})
```

### Error Types
| Error | Handling |
|-------|----------|
| `ErrTokenExpired` | Auto refresh or re-login |
| `ErrInvalidToken` | Force re-login |
| `ErrRateLimited` | Exponential backoff |
| `ErrLoginBlocked` | Surface clear error |

---

## Core Messaging

### Go SDK
```go
// Send
bot.SendText(to, "Hello!")
bot.SendSticker(to, packageID, stickerID)
bot.SendImage(to, imageURL)

// Receive via long polling
bot.OnMessage(func(msg *Message) {
    fmt.Println(msg.From, msg.Text)
})
bot.Listen() // blocking, auto-reconnect
```

**Supported message types:** Text, Sticker, Image, Video, Audio, Location, Contact

---

## OpenChat Features

### Go SDK
```go
oc := bot.OpenChat()

// Discovery & membership
oc.Search("golang developers")
oc.Join(chatID)
oc.Leave(chatID)

// Messaging
oc.SendText(chatID, "Hello!")
oc.OnMessage(chatID, func(msg *OpenChatMessage) { ... })

// Member management
oc.GetMembers(chatID)
oc.Kick(chatID, memberMID)
oc.Ban(chatID, memberMID)

// Admin
oc.PinMessage(chatID, messageID)
oc.SetNotification(chatID, enabled)
oc.GetChatInfo(chatID)
```

---

## HTTP/WebSocket Server

### REST Endpoints
```
POST   /api/message/send
POST   /api/message/send-sticker
GET    /api/contacts
GET    /api/profile/:mid
GET    /api/openchat/search?q=
POST   /api/openchat/join
POST   /api/openchat/leave
POST   /api/openchat/:id/send
GET    /api/openchat/:id/members
DELETE /api/openchat/:id/members/:mid
```

### WebSocket
```
WS /ws/events
← { "type": "message",          "data": { ... } }
← { "type": "openchat_message", "data": { ... } }
← { "type": "auth_error",       "data": { ... } }
```

### Example (start server)
```go
srv := server.New(bot, server.Config{Port: 8080})
srv.Start()
```

---

## Claude Skill: `update-line-api`

**Location:** `skills/update-line-api.md` (repo) + deployed to `~/.claude/plugins/`

**Triggers:**
- `/update-line-api` command
- User pastes error log or HTTP traffic dump

### Three Capabilities

**1. Error Detection & Fix**
- Input: error log / stack trace
- Analyze: identify changed endpoint or parameter
- Output: patch `api/endpoints.go` or `api/talk.go` / `api/openchat.go`

**2. Reverse Engineering Helper**
- Input: HTTP traffic dump (mitmproxy/Charles) or protobuf schema
- Analyze: extract endpoints, headers, payload structure
- Output: update `api/` layer, suggest struct changes

**3. Code + Docs Updater**
- After fix: update `api/` files → update `CHANGELOG.md` → update README API reference → suggest version bump

### Skill Workflow
```
1. Read api/endpoints.go (current state)
2. Analyze input (error log / traffic dump)
3. Identify what changed
4. Edit api/ files
5. Run tests (go test ./...)
6. Update docs
```

---

## Key Design Decisions

| Decision | Choice | Reason |
|----------|--------|--------|
| Architecture | Layered SDK | `api/` layer isolates LINE API changes; easy for Claude skill to target |
| Auth storage | JSON config file | Simple, portable, easy to inspect/edit |
| Server protocol | REST + WebSocket | REST for commands, WebSocket for real-time events |
| OpenChat name | OpenChat (not LINE Square) | Reflects current LINE branding |
| Skill location | In-repo + deployed | Version-controlled alongside API definitions |

---

## Out of Scope (v1)

- LINE Timeline
- LINE Shop / sticker shop
- Voice/video calls
- Multi-account management
