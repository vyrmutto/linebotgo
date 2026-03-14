// Package main is a simple LINE echo bot that replies to every message.
// It uses email/password auth with token persistence.
//
// Usage:
//
//	LINE_EMAIL=xxx@gmail.com LINE_PASSWORD=xxx go run ./examples/simple-bot
package main

import (
	"context"
	"log"
	"os"

	"github.com/vysina/linebotgo/api"
	"github.com/vysina/linebotgo/auth"
	"github.com/vysina/linebotgo/client"
	"github.com/vysina/linebotgo/config"
	"github.com/vysina/linebotgo/transport"
)

const configPath = "linebotgo-config.json"

func main() {
	tok := loadOrLogin()

	bot := client.New(client.WithToken(tok))
	bot.OnMessage(func(msg *client.Message) {
		log.Printf("[%s → %s]: %s", msg.From, msg.To, msg.Text)
		if err := bot.SendText(msg.From, msg.Text); err != nil {
			log.Printf("send error: %v", err)
		}
	})

	log.Println("Echo bot listening...")
	if err := bot.Listen(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func loadOrLogin() *auth.Token {
	cfg, err := config.Load(configPath)
	if err == nil && !isExpired(cfg) {
		log.Println("Loaded token from config")
		return &auth.Token{
			AuthToken:    cfg.AuthToken,
			RefreshToken: cfg.RefreshToken,
			ExpiresAt:    cfg.ExpiresAt,
		}
	}

	email := os.Getenv("LINE_EMAIL")
	password := os.Getenv("LINE_PASSWORD")
	if email == "" || password == "" {
		log.Fatal("Set LINE_EMAIL and LINE_PASSWORD environment variables")
	}

	c := transport.NewHTTPClient()
	tok, err := auth.EmailLogin(c, api.EndpointEmailLogin, email, password)
	if err != nil {
		log.Fatalf("login failed: %v", err)
	}

	if err := config.Save(configPath, &config.Config{
		AuthToken:    tok.AuthToken,
		RefreshToken: tok.RefreshToken,
		ExpiresAt:    tok.ExpiresAt,
	}); err != nil {
		log.Printf("warning: could not save config: %v", err)
	}

	log.Println("Logged in and saved token")
	return tok
}

func isExpired(cfg *config.Config) bool {
	tok := &auth.Token{ExpiresAt: cfg.ExpiresAt}
	return tok.IsExpired()
}
