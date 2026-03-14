// Package main runs linebotgo as an HTTP/WebSocket API server.
// Other services can send messages via REST and receive events via WebSocket.
//
// Usage:
//
//	LINE_EMAIL=xxx@gmail.com LINE_PASSWORD=xxx go run ./examples/api-server
//
// Endpoints:
//
//	POST /api/message/send  {"to":"MID","text":"hello"}
//	GET  /api/contacts
//	WS   /ws/events         <- real-time LINE events
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

const configPath = "linebotgo-config.json"

func main() {
	tok := loadOrLogin()
	bot := client.New(client.WithToken(tok))

	srv := server.New(bot, server.Config{Port: 8080})

	// Forward incoming LINE messages to all WebSocket clients
	bot.OnMessage(func(msg *client.Message) {
		srv.BroadcastMessage(msg.From, msg.To, msg.Text)
	})

	go func() {
		log.Println("API server listening on :8080")
		log.Fatal(srv.Start())
	}()

	log.Println("Long polling LINE events...")
	if err := bot.Listen(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func loadOrLogin() *auth.Token {
	cfg, err := config.Load(configPath)
	if err == nil {
		tok := &auth.Token{ExpiresAt: cfg.ExpiresAt}
		if !tok.IsExpired() {
			log.Println("Loaded token from config")
			return &auth.Token{
				AuthToken:    cfg.AuthToken,
				RefreshToken: cfg.RefreshToken,
				ExpiresAt:    cfg.ExpiresAt,
			}
		}
	}

	email := os.Getenv("LINE_EMAIL")
	password := os.Getenv("LINE_PASSWORD")
	if email == "" || password == "" {
		log.Fatal("Set LINE_EMAIL and LINE_PASSWORD environment variables")
	}

	c := transport.NewHTTPClient()
	tok, err := auth.EmailLogin(c, api.EndpointLogin, email, password)
	if err != nil {
		log.Fatalf("login failed: %v", err)
	}

	config.Save(configPath, &config.Config{ //nolint:errcheck
		AuthToken:    tok.AuthToken,
		RefreshToken: tok.RefreshToken,
		ExpiresAt:    tok.ExpiresAt,
	})

	return tok
}
