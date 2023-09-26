package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config is the configuration for the application
type Config struct {
	BaseURL string `toml:"base_url,omitempty"`
	APIKey  string `toml:"api_key"`
}

// Get reads the configuration from a TOML file and returns a Config
func Get(path string) (Config, error) {
	var (
		cfg Config
		err error
	)

	if path == "" {
		path, err = getDefaultPath()
		if err != nil {
			return Config{}, err
		}
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Write writes the configuration to a TOML file
func Write(path string) error {
	var err error

	cfg := Config{
		APIKey: "CHANGE_ME",
	}

	if path == "" {
		path, err = getDefaultPath()
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(cfg)
}

func getDefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("config err: %w", err)
	}

	return filepath.Join(home, ".certsio.toml"), nil
}
