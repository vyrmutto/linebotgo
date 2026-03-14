package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vysina/linebotgo/config"
)

func TestSaveAndLoadToken(t *testing.T) {
	f, err := os.CreateTemp("", "linebotgo-test-*.json")
	assert.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())

	cfg := &config.Config{
		AuthToken:    "test-token",
		RefreshToken: "test-refresh",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		DeviceID:     "device-123",
	}

	err = config.Save(f.Name(), cfg)
	assert.NoError(t, err)

	loaded, err := config.Load(f.Name())
	assert.NoError(t, err)
	assert.Equal(t, "test-token", loaded.AuthToken)
	assert.Equal(t, "device-123", loaded.DeviceID)
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/tmp/does-not-exist-linebotgo-xyz.json")
	assert.ErrorIs(t, err, config.ErrNoConfig)
}
