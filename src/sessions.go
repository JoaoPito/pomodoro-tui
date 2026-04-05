package main

import (
	"fmt"
	"time"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

func sessionsListWidth(totalWidth int) int {
	h, _ := docStyle.GetFrameSize()
	return (totalWidth - h) * 60 / 100
}

func (m model) startSessionsView() (model, tea.Cmd) {
	listItems := Map(m.availSessions, func(s sessionType) list.Item {
		return item{title: s.Name, desc: s.Duration.String()}
	})

	m.state = sessionsView
	m.list.SetItems(listItems)

	completedTitle := ""
	if m.tasks[m.selTask].Completed {
		completedTitle = "[COMPLETED] "
	}

	m.list.Title = "> " + completedTitle + m.tasks[m.selTask].Name + " - Choose a focus session"

	if m.width > 0 {
		_, v := docStyle.GetFrameSize()
		m.list.SetSize(sessionsListWidth(m.width), m.height-v)
	}

	return m, nil
}

func (m model) updateSessions(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "enter":
			m.selSessionType = m.list.Index()
			return m.startTimerView()

		case "backspace", "b":
			return m.startTasksView()

		case "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		_, v := docStyle.GetFrameSize()
		m.list.SetSize(sessionsListWidth(msg.Width), msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) viewSessions() string {
	h, _ := docStyle.GetFrameSize()
	availWidth := m.width - h
	listWidth := sessionsListWidth(m.width)
	descWidth := availWidth - listWidth

	desc := m.tasks[m.selTask].Description
	if desc == "" {
		desc = "(no description)"
	}
	if m.tasks[m.selTask].EstimatedDuration != nil {
		estimatedDurationMins := (*m.tasks[m.selTask].EstimatedDuration) * time.Minute 
		desc += lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf("\n\nEstim. duration: %s", estimatedDurationMins.String()))
	}

	descPanel := lipgloss.NewStyle().
		Width(descWidth - 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lightGray).
		Padding(1, 1).
		Render(desc)

	return lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), descPanel)
}
