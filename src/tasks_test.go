package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	api "pomodoro-tui/apiclient"
)

func newTasksTestModel(apiURL string, tasks []api.Task) model {
	client := api.NewClient(apiURL, "test-key", "test-device")
	m := model{
		apiclient:  client,
		state:      tasksView,
		tasks:      tasks,
		projects:   []api.Project{{ID: 1, Name: "Test Project"}},
		selProject: 0,
		list:       list.New(nil, list.NewDefaultDelegate(), 80, 24),
	}
	items := Map(tasks, mapTaskToListItem)
	m.list.SetItems(items)
	return m
}

func TestUpdateTasks_PressDToOpenConfirmation(t *testing.T) {
	tasks := []api.Task{{ID: 1, Name: "Buy groceries"}}
	m := newTasksTestModel("http://localhost", tasks)

	result, _ := m.updateTasks(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	resultModel := result.(model)

	if resultModel.state != deleteTaskView {
		t.Errorf("expected state=%d (deleteTaskView), got %d", deleteTaskView, resultModel.state)
	}
}

func TestUpdateTasks_PressDOnEmptyListDoesNothing(t *testing.T) {
	m := newTasksTestModel("http://localhost", []api.Task{})

	result, _ := m.updateTasks(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	resultModel := result.(model)

	if resultModel.state != tasksView {
		t.Errorf("expected state to remain tasksView (%d), got %d", tasksView, resultModel.state)
	}
}

func TestUpdateDeleteTask_PressNToCancelDeletion(t *testing.T) {
	tasks := []api.Task{{ID: 1, Name: "Buy groceries"}}
	m := newTasksTestModel("http://localhost", tasks)
	m.state = deleteTaskView

	result, _ := m.updateDeleteTask(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	resultModel := result.(model)

	if resultModel.state != tasksView {
		t.Errorf("expected state=%d (tasksView), got %d", tasksView, resultModel.state)
	}
}

func TestUpdateDeleteTask_PressEscToCancelDeletion(t *testing.T) {
	tasks := []api.Task{{ID: 1, Name: "Some task"}}
	m := newTasksTestModel("http://localhost", tasks)
	m.state = deleteTaskView

	result, _ := m.updateDeleteTask(tea.KeyMsg{Type: tea.KeyEsc})
	resultModel := result.(model)

	if resultModel.state != tasksView {
		t.Errorf("expected state=%d (tasksView), got %d", tasksView, resultModel.state)
	}
}

func TestUpdateDeleteTask_PressYToConfirmDeletion(t *testing.T) {
	deleteCallCount := 0
	tasks := []api.Task{{ID: 42, Name: "Buy groceries"}}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/tasks" && r.Method == http.MethodDelete:
			deleteCallCount++
			json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "id": 42})
		case r.URL.Path == "/tasks/get-by-project":
			json.NewEncoder(w).Encode(tasks)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	m := newTasksTestModel(srv.URL, tasks)
	m.state = deleteTaskView
	m.selTask = 0

	result, _ := m.updateDeleteTask(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	resultModel := result.(model)

	if deleteCallCount != 1 {
		t.Errorf("expected DeleteTask to be called once, got %d", deleteCallCount)
	}
	if resultModel.state != tasksView {
		t.Errorf("expected state=%d (tasksView), got %d", tasksView, resultModel.state)
	}
}
