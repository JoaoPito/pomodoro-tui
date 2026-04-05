package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type SessionTypeConfig struct {
	Name            string `json:"name"`
	DurationMinutes int    `json:"duration_minutes"`
}

func (s SessionTypeConfig) Duration() time.Duration {
	return time.Duration(s.DurationMinutes) * time.Minute
}

type Config struct {
	APIUrl       string              `json:"api_url"`
	APIKey       string              `json:"api_key"`
	DeviceName   string              `json:"device_name"`
	SessionTypes []SessionTypeConfig `json:"session_types"`
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("could not read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("could not parse config file: %w", err)
	}

	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func validate(cfg Config) error {
	if cfg.APIUrl == "" {
		return errors.New("config: api_url is required")
	}
	if cfg.APIKey == "" {
		return errors.New("config: api_key is required")
	}
	if cfg.DeviceName == "" {
		return errors.New("config: device_name is required")
	}
	if len(cfg.SessionTypes) == 0 {
		return errors.New("config: session_types must not be empty")
	}
	for i, s := range cfg.SessionTypes {
		if s.Name == "" {
			return fmt.Errorf("config: session_types[%d]: name is required", i)
		}
		if s.DurationMinutes <= 0 {
			return fmt.Errorf("config: session_types[%d]: duration_minutes must be greater than 0", i)
		}
	}
	return nil
}
