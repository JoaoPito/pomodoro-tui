package main

import (
	api "pomodoro-tui/apiclient"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/timer"
	"time"
)

const (
	projectsView uint = iota
	tasksView
	sessionsView
	timerView
	newTaskView
)

type sessionType struct { 
	Name 		string
	Duration 	time.Duration
}

type model struct {
	state 		uint
	apiclient 	*api.Client

	projects 	[]api.Project
	tasks 		[]api.Task

	selProject 	int
	selTask 	int	
	selSessionType 	int	
	session 	*api.FocusSession

	list 		list.Model
	timer		timer.Model
	newTaskForm	newTaskForm	

	availSessions 	[]sessionType

	width 		int
	height 		int
}

func NewModel(apiClient *api.Client, sessions []sessionType) model {
	m := model{
		apiclient: apiClient,
		state: projectsView,
		projects: nil,
		tasks: nil,
		selProject: 0,
		selTask: 0,
		selSessionType: 0,
		session: nil,
		list: list.New(nil, list.NewDefaultDelegate(), 0, 0),
		timer: timer.New(0),
		newTaskForm: CreateNewTaskForm(),
		availSessions: sessions,
	}
	m, _ = m.startProjectsView()
	return m
}

func (m model) Init() tea.Cmd {
	return nil 
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sz, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = sz.Width
		m.height = sz.Height
	}
	switch m.state {
	case projectsView:
		return m.updateProjects(msg)
	case tasksView:
		return m.updateTasks(msg)
	case sessionsView:
		return m.updateSessions(msg)
	case timerView:
		return m.updateTimer(msg)
	case newTaskView:
		return m.updateNewTask(msg)
	}

	return m, nil
}
