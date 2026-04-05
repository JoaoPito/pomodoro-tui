package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/textarea"
	"time"
	"fmt"
)

const(
	name = iota
	prior
	deadl
	estim
)

type newTaskForm struct {
	inputs	[]textinput.Model
	description textarea.Model
	focused	int
	err	error
}

func (m *newTaskForm) NextInput() {
	m.focused = (m.focused + 1) % (len(m.inputs) + 1)
}

func (m *newTaskForm) PrevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs)
	}
}

func (m *newTaskForm) FocusOnSelected(){

	for i := range m.inputs {
		m.inputs[i].Blur()
	}
	m.description.Blur()

	if m.focused == len(m.inputs) {
		m.description.Focus()
	} else {
		m.inputs[m.focused].Focus()
	}
}

func nameValidator(s string) error {
	if len(s) <= 0 {
		return fmt.Errorf("Name cannot be empty!")
	}
	return nil
}

func priorityValidator(s string) error {
	if !(s == "" || s == "!" || s == "!!" || s == "!!!" ){
		return fmt.Errorf("Priority invalid!")
	}
	return nil
}

func estimatedDurationValidator(s string) error {
	if s == "" {
		return nil
	}

	_, err := time.ParseDuration(s)
	return err
}

func deadlineValidator(s string) error {
	if s == "" {
		return nil
	}
	_, err := time.Parse("02-01", s)
	return err
}

func CreateNewTaskForm() newTaskForm {
	var form newTaskForm

	var inputs []textinput.Model = make([]textinput.Model, 4)
	form.inputs = inputs
	form.err = nil

	form.inputs[name] = textinput.New()
	form.inputs[name].Placeholder = "Buy coffee"
	form.inputs[name].Focus()
	form.inputs[name].Width = 80
	form.inputs[name].Prompt = ""
	form.inputs[name].Validate = nameValidator

	form.inputs[prior] = textinput.New()
	form.inputs[prior].Placeholder = "!!!"
	form.inputs[prior].Width = 3
	form.inputs[prior].CharLimit = 3
	form.inputs[prior].Prompt = ""
	form.inputs[prior].Validate = priorityValidator

	form.inputs[estim] = textinput.New()
	form.inputs[estim].Placeholder = "1h30m"
	form.inputs[estim].Width = 10
	form.inputs[estim].CharLimit = 10
	form.inputs[estim].Prompt = ""
	form.inputs[estim].Validate = estimatedDurationValidator
	
	form.inputs[deadl] = textinput.New()
	form.inputs[deadl].Placeholder = "DD-MM"
	form.inputs[deadl].Width = 5
	form.inputs[deadl].CharLimit = 5
	form.inputs[deadl].Prompt = ""
	form.inputs[deadl].Validate = deadlineValidator

	form.description = textarea.New()
	form.description.Placeholder = "Just a castaway, an island lost at sea, oh. Another lonely day with no one here but me, oh"
	form.description.Prompt = "┃ "
	form.description.CharLimit = 280
	form.description.SetWidth(80)
	form.description.SetHeight(5)
	form.description.ShowLineNumbers = true

	return form
}
