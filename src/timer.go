package main

import (
	"log"
	"time"

	api "pomodoro-tui/apiclient"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/gen2brain/beeep"
)

func (m model) startTimerView() (model, tea.Cmd) {
	m.state = timerView
	
	// POST session API
	session := api.FocusSession {
		ID: nil,
		Start: time.Now(),
		End: nil,
		TaskID: m.tasks[m.selTask].ID,
		Session: m.availSessions[m.selSessionType].Name,
	}
	sessionId, err := m.apiclient.UpsertFocusSession(session)

	if err != nil {
		log.Fatalf("failed to send session data to API: %v", err)
	}

	// save session info on model
	session.ID = &sessionId
	m.session = &session
	
	m.timer = timer.New(m.availSessions[m.selSessionType].Duration)
	beeep.AppName = "pomodoro-tui"
		
	return m, tea.Batch(m.timer.Init(), m.spinner.Tick)
}

func (m model) updateTimer(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	
	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		//m.keymap.stop.SetEnabled(m.timer.Running())
		//m.keymap.start.SetEnabled(!m.timer.Running())
		//m = m.stopTimer()
		return m, cmd

	case timer.TimeoutMsg:
		return m.stopTimer()

	case tea.KeyMsg:
		switch key := msg.String(); key{
		case "backspace", "b":
			return m.stopTimer()
			// Save pomodoro session with API
		case "ctrl+c":
			m.stopTimer()
			return m, tea.Quit
		}
	}

	return m, nil 
}

func (m model) viewTimer() string {
	task := m.tasks[m.selTask]

	s := titleStyle.Render(boldTitleStyle.Render("> " + task.Name)) + "\n\n"
	s += subTextStyle.Render(task.Description) + "\n\n"
	s += timerStyle.Render(spinnerStyle.Render(m.spinner.View() + " ") + m.timer.View() + " remaining...")
	
	if m.timer.Timedout() {
		s = "All done!\n"
	}

	return s
}

func truncateSessionEnd(session *api.FocusSession, duration time.Duration) {
	if session.End == nil {
		return
	}
	if session.End.Sub(session.Start) > duration {
		truncated := session.Start.Add(duration)
		session.End = &truncated
	}
}

func (m model) stopTimer() (model, tea.Cmd) {
	m.timer.Stop()

	now := time.Now()
	m.session.End = &now
	truncateSessionEnd(m.session, m.availSessions[m.selSessionType].Duration)
	_, err := m.apiclient.UpsertFocusSession(*m.session)
	if err != nil {
		log.Fatalf("failed to send session data to API: %v", err)
	}

	err = beeep.Notify(m.session.Session + " session finished!", "nice work!", "")
	if err != nil {
        	panic(err)
    	}

	err = beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)

	m.session = nil
	m.state = sessionsView
	return m, nil
}

func CreateNewSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Jump
	return s
}
