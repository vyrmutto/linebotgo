package server_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vysina/linebotgo/server"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func TestWebSocketBroadcast(t *testing.T) {
	srv := server.New(&mockSender{}, server.Config{})
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	wsURL := "ws" + ts.URL[4:] + "/ws/events"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close(websocket.StatusNormalClosure, "done")

	// Small delay to ensure connection is registered in hub
	time.Sleep(100 * time.Millisecond)

	// Broadcast from server
	srv.Broadcast(server.Event{Type: "message", Data: map[string]string{"text": "hello ws"}})

	var received server.Event
	err = wsjson.Read(ctx, conn, &received)
	assert.NoError(t, err)
	assert.Equal(t, "message", received.Type)
}
