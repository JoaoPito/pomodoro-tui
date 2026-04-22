package apiclient_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pomodoro-tui/apiclient"
)

func TestDeleteTask_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/tasks" {
			t.Errorf("expected path /tasks, got %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if id, ok := body["id"].(float64); !ok || uint(id) != 42 {
			t.Errorf("expected body id=42, got %v", body["id"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "id": 42})
	}))
	defer srv.Close()

	client := apiclient.NewClient(srv.URL, "test-key", "test-device")
	err := client.DeleteTask(42)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestDeleteTask_APIError(t *testing.T) {
	errMsg := "task not found"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   errMsg,
		})
	}))
	defer srv.Close()

	client := apiclient.NewClient(srv.URL, "test-key", "test-device")
	err := client.DeleteTask(99)
	if err == nil {
		t.Fatal("expected error for API failure, got nil")
	}
}

func TestDeleteTask_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := apiclient.NewClient(srv.URL, "test-key", "test-device")
	err := client.DeleteTask(1)
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}
}
