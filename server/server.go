package server

import (
	"fmt"
	"net/http"
)

// Sender is the interface the server uses to send messages.
// Decouples server from the concrete client implementation.
type Sender interface {
	SendText(to, text string) error
}

// Config holds server configuration.
type Config struct {
	Port int
}

// Server serves LINE bot operations over REST and WebSocket.
type Server struct {
	sender Sender
	cfg    Config
	hub    *wsHub
	mux    *http.ServeMux
}

// New creates a Server wrapping the given Sender.
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

// Handler returns the HTTP handler (used in tests without starting a server).
func (s *Server) Handler() http.Handler {
	return s.mux
}

// Start starts the HTTP server on the configured port.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.cfg.Port)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/api/message/send", s.handleSendMessage)
	s.mux.HandleFunc("/api/contacts", s.handleGetContacts)
	s.mux.HandleFunc("/ws/events", s.handleWebSocket)
}
