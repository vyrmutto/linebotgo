package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

var ErrNoConfig = errors.New("config file not found")

// Config holds persisted LINE auth state between sessions.
type Config struct {
	AuthToken    string    `json:"auth_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	DeviceID     string    `json:"device_id"`
	Certificate  string    `json:"certificate"`
}

// Save writes cfg as JSON to path (mode 0600 — owner read/write only).
func Save(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// Load reads Config from path. Returns ErrNoConfig if file missing.
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
