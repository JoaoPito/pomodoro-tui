package main

import (
	"log"

	api "pomodoro-tui/apiclient"
	"pomodoro-tui/config"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
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
