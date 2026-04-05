package main

import (
	"log"
	"time"
	api "pomodoro-tui/apiclient"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	client := api.NewClient(
		"https://cacambitos.joaod.com/webhook/projects",
		"12345",
		"archlinux-btw",
	)

	sessions := []sessionType {
		{"work", 45 * time.Minute},
		{"short work", 25 * time.Minute},
		{"break", 5 * time.Minute},
		{"long break", 15 * time.Minute},
	}

	m := NewModel(client, sessions)
	
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("unable to run tui: %v", err)
	}
}
