package server

import (
	"context"
	"net/http"
	"sync"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Event is a real-time event broadcast to all WebSocket clients.
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
	defer h.mu.Unlock()
	h.conns[c] = struct{}{}
}

func (h *wsHub) remove(c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conns, c)
}

func (h *wsHub) broadcast(e Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.conns {
		// best-effort; ignore send errors for individual clients
		wsjson.Write(context.Background(), c, e) //nolint:errcheck
	}
}

// Broadcast sends an event to all connected WebSocket clients.
func (s *Server) Broadcast(e Event) {
	s.hub.broadcast(e)
}

// BroadcastMessage is a convenience helper called when a LINE message arrives.
func (s *Server) BroadcastMessage(from, to, text string) {
	s.Broadcast(Event{
		Type: "message",
		Data: map[string]string{"from": from, "to": to, "text": text},
	})
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}
	s.hub.add(conn)
	defer s.hub.remove(conn)

	// CloseRead drains incoming frames; returns when client disconnects.
	ctx := conn.CloseRead(context.Background())
	<-ctx.Done()
}
