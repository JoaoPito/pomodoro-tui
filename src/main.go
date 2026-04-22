package main

import (
	"log"
	"os"
	"path/filepath"

	api "pomodoro-tui/apiclient"
	"pomodoro-tui/config"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("unable to resolve executable path: %v", err)
	}
	configPath := filepath.Join(filepath.Dir(exe), "config.json")

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	client := api.NewClient(cfg.APIUrl, cfg.APIKey, cfg.DeviceName)

	sessions := make([]sessionType, len(cfg.SessionTypes))
	for i, s := range cfg.SessionTypes {
		sessions[i] = sessionType{s.Name, s.Duration()}
	}

	m := NewModel(client, sessions)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("unable to run tui: %v", err)
	}
}
