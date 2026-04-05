package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"pomodoro-tui/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}

func TestLoadConfig_Valid(t *testing.T) {
	path := writeTempConfig(t, `{
		"api_url": "https://example.com/webhook",
		"api_key": "secret123",
		"device_name": "my-machine",
		"session_types": [
			{"name": "work", "duration_minutes": 45},
			{"name": "break", "duration_minutes": 5}
		]
	}`)

	cfg, err := config.LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.APIUrl != "https://example.com/webhook" {
		t.Errorf("APIUrl = %q, want %q", cfg.APIUrl, "https://example.com/webhook")
	}
	if cfg.APIKey != "secret123" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "secret123")
	}
	if cfg.DeviceName != "my-machine" {
		t.Errorf("DeviceName = %q, want %q", cfg.DeviceName, "my-machine")
	}
	if len(cfg.SessionTypes) != 2 {
		t.Fatalf("len(SessionTypes) = %d, want 2", len(cfg.SessionTypes))
	}
	if cfg.SessionTypes[0].Name != "work" {
		t.Errorf("SessionTypes[0].Name = %q, want %q", cfg.SessionTypes[0].Name, "work")
	}
	if cfg.SessionTypes[0].Duration() != 45*time.Minute {
		t.Errorf("SessionTypes[0].Duration() = %v, want %v", cfg.SessionTypes[0].Duration(), 45*time.Minute)
	}
	if cfg.SessionTypes[1].Name != "break" {
		t.Errorf("SessionTypes[1].Name = %q, want %q", cfg.SessionTypes[1].Name, "break")
	}
	if cfg.SessionTypes[1].Duration() != 5*time.Minute {
		t.Errorf("SessionTypes[1].Duration() = %v, want %v", cfg.SessionTypes[1].Duration(), 5*time.Minute)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := config.LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadConfig_MalformedJSON(t *testing.T) {
	path := writeTempConfig(t, `{ "api_url": "https://example.com", invalid }`)

	_, err := config.LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

func TestLoadConfig_MissingAPIUrl(t *testing.T) {
	path := writeTempConfig(t, `{
		"api_key": "secret",
		"device_name": "machine",
		"session_types": [{"name": "work", "duration_minutes": 25}]
	}`)

	_, err := config.LoadConfig(path)
	if err == nil {
		t.Fatal("expected validation error for missing api_url, got nil")
	}
}

func TestLoadConfig_MissingAPIKey(t *testing.T) {
	path := writeTempConfig(t, `{
		"api_url": "https://example.com",
		"device_name": "machine",
		"session_types": [{"name": "work", "duration_minutes": 25}]
	}`)

	_, err := config.LoadConfig(path)
	if err == nil {
		t.Fatal("expected validation error for missing api_key, got nil")
	}
}

func TestLoadConfig_MissingDeviceName(t *testing.T) {
	path := writeTempConfig(t, `{
		"api_url": "https://example.com",
		"api_key": "secret",
		"session_types": [{"name": "work", "duration_minutes": 25}]
	}`)

	_, err := config.LoadConfig(path)
	if err == nil {
		t.Fatal("expected validation error for missing device_name, got nil")
	}
}

func TestLoadConfig_EmptySessionTypes(t *testing.T) {
	path := writeTempConfig(t, `{
		"api_url": "https://example.com",
		"api_key": "secret",
		"device_name": "machine",
		"session_types": []
	}`)

	_, err := config.LoadConfig(path)
	if err == nil {
		t.Fatal("expected validation error for empty session_types, got nil")
	}
}

func TestLoadConfig_SessionTypeMissingName(t *testing.T) {
	path := writeTempConfig(t, `{
		"api_url": "https://example.com",
		"api_key": "secret",
		"device_name": "machine",
		"session_types": [{"duration_minutes": 25}]
	}`)

	_, err := config.LoadConfig(path)
	if err == nil {
		t.Fatal("expected validation error for session type missing name, got nil")
	}
}

func TestLoadConfig_SessionTypeZeroDuration(t *testing.T) {
	path := writeTempConfig(t, `{
		"api_url": "https://example.com",
		"api_key": "secret",
		"device_name": "machine",
		"session_types": [{"name": "work", "duration_minutes": 0}]
	}`)

	_, err := config.LoadConfig(path)
	if err == nil {
		t.Fatal("expected validation error for session type with zero duration, got nil")
	}
}
