package main

import (
	"time"
	"log"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"fmt"
	"strings"
	api "pomodoro-tui/apiclient"
)

type (
	errMsg error
)

func (m model) startNewTaskView() (tea.Model, tea.Cmd) {
	m.state = newTaskView

	for _,input := range m.newTaskForm.inputs {
		input.Reset()
	}

	return m, textinput.Blink
}

func (m model) updateNewTask(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key{
		case "alt+enter":
			if m.validateTaskForm() != nil {
				return m, nil
			}
			// Create task
			task := m.formToTask()
			err := m.apiclient.InsertTask(task)
			if err != nil {
				log.Fatalf("error creating task: %v", err)
			}
			// Call api
			return m.startTasksView()

		case "esc":
			// Go to tasks view without saving
			return m.startTasksView()

		case "shift+tab", "ctrl+p":
			m.newTaskForm.PrevInput()

		case "tab", "ctrl+n":
			m.newTaskForm.NextInput()

		case "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case errMsg:
		m.newTaskForm.err = msg
		return m, nil
	
	}

	for i := range m.newTaskForm.inputs {
		var cmd tea.Cmd
		m.newTaskForm.inputs[i], cmd = m.newTaskForm.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	
	// description text area
	var descCmd tea.Cmd
	m.newTaskForm.description, descCmd = m.newTaskForm.description.Update(msg)
	cmds = append(cmds, descCmd)

	m.newTaskForm.FocusOnSelected()

	return m, tea.Batch(cmds...)
}

func (m model) validateTaskForm() error {
	if err := nameValidator(m.newTaskForm.inputs[name].Value()); err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	if err := estimatedDurationValidator(m.newTaskForm.inputs[estim].Value()); err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	if err := deadlineValidator(m.newTaskForm.inputs[deadl].Value()); err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	if err := priorityValidator(m.newTaskForm.inputs[prior].Value()); err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	return nil
}

func (m model) formToTask() api.Task {
	f := m.newTaskForm

	var priority uint
	priority = uint(strings.Count(f.inputs[prior].Value(), "!"))

	deadlineInput := f.inputs[deadl].Value()
	var deadline *time.Time
	if deadlineInput != "" {
		year := time.Now().Year()
		deadlineTime, _ := time.Parse("02-01-2006", fmt.Sprintf("%s-%d", deadlineInput, year))
		deadline = &deadlineTime
	}

	var estimDuration *time.Duration
	estimDurationInput := f.inputs[estim].Value()
	if estimDurationInput != "" {
		estimDurationParsed, _ := time.ParseDuration(estimDurationInput)
		estimDuration = &estimDurationParsed
	}

	return api.Task{
		ProjectID: m.projects[m.selProject].ID,
		Name: f.inputs[name].Value(),
		Completed: false,
		Priority: priority,
		Description: f.description.Value(),
		Deadline: deadline,
		EstimatedDuration: estimDuration,
	}
}

func (m model) viewNewTask() string {
	projectName := m.projects[m.selProject].Name
	var errorMsg string

	if m.newTaskForm.err != nil {
		errorMsg = m.newTaskForm.err.Error()
	}

	return fmt.Sprintf(
		`%s - New task


%s
%s

%s  %s  %s
%s		%s		%s

%s
%s


%s
%s
`,
		projectName,
		inputStyle.Width(30).Render("Name"),
		m.newTaskForm.inputs[name].View(),
		inputStyle.Width(10).Render("Priority"),
		inputStyle.Width(10).Render("Deadline"),
		inputStyle.Width(15).Render("Estim. Time"),
		m.newTaskForm.inputs[prior].View(),
		m.newTaskForm.inputs[deadl].View(),
		m.newTaskForm.inputs[estim].View(),
		inputStyle.Width(50).Render("Description"),
		m.newTaskForm.description.View(),
		errorMsg,
		continueStyle.Render("ALT + ENTER to save and continue"),
	) + "\n"
}
