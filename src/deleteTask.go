package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) startDeleteTaskView() (model, tea.Cmd) {
	m.state = deleteTaskView
	return m, nil
}

func (m model) updateDeleteTask(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			taskID := m.tasks[m.selTask].ID
			if err := m.apiclient.DeleteTask(taskID); err != nil {
				log.Fatalf("failed to delete task: %v", err)
			}
			return m.startTasksView()

		case "n", "esc", "q":
			m.state = tasksView
			return m, nil

		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) viewDeleteTask() string {
	taskName := m.tasks[m.selTask].Name

	content := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		deleteConfirmStyle.Render("> Delete task?"),
		taskName,
		deleteHintStyle.Render("[y] Yes    [n] No"),
	)

	popup := deletePopupStyle.Render(content)

	h, v := docStyle.GetFrameSize()
	return lipgloss.Place(
		m.width-h,
		m.height-v,
		lipgloss.Center,
		lipgloss.Center,
		popup,
	)
}
