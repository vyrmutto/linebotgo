package client_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vysina/linebotgo/auth"
	"github.com/vysina/linebotgo/client"
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

func TestNewClientNoToken(t *testing.T) {
	c := client.New()
	assert.Equal(t, "", c.AuthToken())
}
