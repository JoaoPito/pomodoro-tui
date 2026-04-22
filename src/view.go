package main

import (
	"github.com/charmbracelet/lipgloss"
)


const (
	hotPink = lipgloss.Color("13")
	complBlue = lipgloss.Color("12")
	darkGray = lipgloss.Color("#767676")
	lightGray = lipgloss.Color("#808080")
	black = lipgloss.Color("232")
)

var (
	docStyle 	= lipgloss.NewStyle().Margin(5, 5)
	inputStyle    	= lipgloss.NewStyle().Foreground(hotPink)
	continueStyle 	= lipgloss.NewStyle().Foreground(darkGray)
	completedStyle 	= lipgloss.NewStyle().Strikethrough(true).Foreground(lightGray)
	timerStyle	= lipgloss.NewStyle().Bold(true).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(hotPink)).Padding(0, 16)
	titleStyle	= lipgloss.NewStyle().Background(hotPink).Foreground(black).Padding(0,1)
	boldTitleStyle 	= lipgloss.NewStyle().Bold(true)
	subTextStyle 	= lipgloss.NewStyle().Foreground(lightGray)
	deletePopupStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(hotPink).
				Padding(1, 4)
	deleteConfirmStyle = lipgloss.NewStyle().Background(hotPink).Foreground(black).Padding(0,1).Bold(true)
	deleteHintStyle    = lipgloss.NewStyle().Foreground(lightGray)
	spinnerStyle	= lipgloss.NewStyle().Foreground(complBlue)
)

func (m model) View() string {
	var s string
	switch m.state {
		case projectsView:
			s = m.viewProjects()
		case tasksView:
			s = m.viewTasks()
		case sessionsView:
			s = m.viewSessions()
		case timerView:
			s = m.viewTimer()
		case newTaskView:
			s = m.viewNewTask()
		case deleteTaskView:
			s = m.viewDeleteTask()
	}

	return docStyle.Render(s)

}
