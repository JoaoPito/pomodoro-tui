package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"slices"
	api "pomodoro-tui/apiclient"
	"github.com/charmbracelet/bubbles/list"
)

func (m model) startProjectsView() (model, tea.Cmd) {
	projects, err := m.apiclient.GetLatestProjects()
	if err != nil {
		log.Fatalf("unable to get latest projects: %v", err)
	}

	slices.SortFunc(projects, func(a, b api.Project) int {
		return int(b.Last_task_updated_at.Unix() - a.Last_task_updated_at.Unix())
	})

	var listItems []list.Item
	listItems = Map(projects, func(p api.Project) list.Item{
		return item { title: p.Name, desc: p.Repository }
	})

	m.projects = projects
	m.list.SetItems(listItems)
	m.state = projectsView
	m.list.Title = "> Choose a project to work on"

	return m, nil
}

func (m model) updateProjects(msg tea.Msg) (tea.Model, tea.Cmd) {	
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key{
		case "enter":
			m.selProject = m.list.Index()
			return m.startTasksView()

		case "ctrl+c":
			return m, tea.Quit
		}
		
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) viewProjects() string {
	return m.list.View()
}
