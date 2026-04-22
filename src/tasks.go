package main

import (
	"time"
	"fmt"
	"log"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/key"
	api "pomodoro-tui/apiclient"
	"sort"
	"github.com/dustin/go-humanize"
)

func (m model) startTasksView() (model, tea.Cmd) {
	m = m.refreshTaskList()
	m.state = tasksView
	m.list.Title = "> " + m.projects[m.selProject].Name
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "add new task"),
			),
			key.NewBinding(
				key.WithKeys("c", "tab"),
				key.WithHelp("c/tab", "toggle completed"),
			),
			key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "delete task"),
			),
		}
	}

	return m, nil
}

func mapTaskToListItem(t api.Task) list.Item {
	var priorityStr string
	if t.Priority != 0 {
		priorityStr = fmt.Sprintf("[%s] ", strings.Repeat("!", int(t.Priority)))
	}
	
	deadlineStr := "no deadline"
	if t.Deadline != nil {
		y, m, d := t.Deadline.Date()
		loc, _ := time.LoadLocation("America/Sao_Paulo")
		deadlineEndOfDay := time.Date(y, m, d, 23, 59, 59, 0, loc)
		deadlineStr = fmt.Sprintf("%s", humanize.Time(deadlineEndOfDay))
	}

	name := fmt.Sprintf("%s%s", priorityStr, t.Name)
	description := deadlineStr
	
	if t.Completed == true {
		name = completedStyle.Render(name)
		description = completedStyle.Render(description)
	}

	return item { title: name, desc: description }
}

func (m model) updateTasks(msg tea.Msg) (tea.Model, tea.Cmd) {	
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key{
		case "enter":
			m.selTask = m.list.Index()
			return m.startSessionsView()

		case "backspace", "b":
			return m.startProjectsView()
		
		case "a":
			return m.startNewTaskView()
		
		case "c", "tab":
			m.toggleTaskCompletion()
			return m.startTasksView()
		
		case "d":
			if len(m.tasks) == 0 {
				return m, nil
			}
			m.selTask = m.list.Index()
			return m.startDeleteTaskView()
		
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

func (m model) viewTasks() string {
	return m.list.View()
}

func (m model) refreshTaskList() model {
	tasks, err := m.apiclient.GetWeekTasksByProject(m.projects[m.selProject].ID)
	if err != nil {
		log.Fatalf("unable to get tasks for project: %v", err)
	}

	sort.Slice(tasks, func(i, j int) bool {
		ti, tj := tasks[i], tasks[j]
		if ti.Completed != tj.Completed {
			return !ti.Completed
		}
		if ti.Priority != tj.Priority {
			return ti.Priority > tj.Priority
		}
		if ti.Deadline == nil && tj.Deadline == nil {
			return false
		}
		if ti.Deadline == nil {
			return false
		}
		if tj.Deadline == nil {
			return true
		}
		return ti.Deadline.Before(*tj.Deadline)
	})

	var listItems []list.Item
	listItems = Map(tasks, mapTaskToListItem)

	m.tasks = tasks
	m.list.SetItems(listItems)
	return m
}

func (m model) toggleTaskCompletion() {
	selectedTask := m.list.Index()
	taskId := m.tasks[selectedTask].ID
	taskCompleted := m.tasks[selectedTask].Completed

	err := m.apiclient.UpdateTaskCompletion(taskId, !taskCompleted)
	if err != nil {
		log.Fatalf("failed to update task completed flag: %v", err)
	}
}
