# linebotgo

A Go library for interacting with LINE's messaging service, built for educational purposes to understand how external tools can interface with LINE.

> **Educational use only.** This project is intended solely for studying how LINE's internal protocol works. It is not affiliated with, endorsed by, or officially supported by LINE Corporation. Use in accordance with LINE's Terms of Service.

## Acknowledgements

Protocol research and reference implementation based on [line-py](https://github.com/fadhiilrachman/line-py) by fadhiilrachman. API constants extracted from the LINE Android APK via [jadx](https://github.com/skylot/jadx) decompilation.

## Features

- Thrift RPC transport over HTTP/2 (LEGY gateway)
- Authentication: email/password and QR code login
- Messaging: send and receive messages, poll for operations
- OpenChat (LINE Square): join, send messages, manage members
- Full endpoint coverage — all 46 service paths from `EnumC55943a`

## Requirements

- Go 1.21+

## Installation

```bash
go get github.com/vyrmutto/linebotgo
```

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/vyrmutto/linebotgo/auth"
    "github.com/vyrmutto/linebotgo/client"
    "github.com/vyrmutto/linebotgo/config"
)

func main() {
    cfg := config.New()

    // Login with email
    token, err := auth.LoginWithEmail(cfg, "your@email.com", "password")
    if err != nil {
        log.Fatal(err)
    }

    c := client.New(cfg, token)

    profile, err := c.GetProfile()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Logged in as:", profile.DisplayName)
}
```

## Project Structure

```
api/         — Endpoint constants and Thrift type definitions
auth/        — Login flows (email, QR code)
client/      — TalkService methods (messages, contacts, poll)
config/      — Configuration
openchat/    — LINE Square / OpenChat service
server/      — Webhook server helper
transport/   — HTTP client with LINE headers
```

## Updating API Constants

When LINE releases a new APK, run the `/update-line-api` skill (see `skills/update-line-api.md`) to re-extract all constants from the decompiled APK.

## Disclaimer

This project is for **educational and research purposes only**. The authors do not condone use of this library to violate LINE's Terms of Service, spam users, or engage in any harmful activity. Use responsibly.

## License

MIT — see [LICENSE](LICENSE)
