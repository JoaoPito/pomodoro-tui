package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"github.com/google/uuid"
	"log"
)

type Client struct {
	BaseURL 	string
	HTTPClient 	*http.Client
	APIKey 		string
	DeviceName	string
}

func NewClient(baseURL, apiKey string, deviceName string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey: apiKey,
		DeviceName: deviceName,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) makeRequest(req *http.Request, v interface{}) error{
	req.Header.Set("Authorization", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return fmt.Errorf("error calling n8n API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("error decoding n8n API response: %w", err)
		}
	}

	return nil
}

func (c *Client) newRequest(method, path string, body interface{}, result interface{}) error {
    url := c.BaseURL + path
    
    var bodyReader io.Reader
    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return fmt.Errorf("erro ao serializar body: %w", err)
        }
        bodyReader = bytes.NewBuffer(jsonData)
    }
    
    req, err := http.NewRequest(method, url, bodyReader)
    if err != nil {
        return fmt.Errorf("erro ao criar request: %w", err)
    }
    
    return c.makeRequest(req, result)
}


func (c *Client) GetLatestProjects() ([]Project, error) {
	var projects []Project
	
	now := time.Now()
//	firstDayOfCurMonth := time.Date(
//		now.Year(),
//		now.Month(),
//		1,
//		0, 0, 0, 0,
//		now.Location(),
//	)

	aYearAgo := now.AddDate(-1, 0, 0)
	
	request := GetLatestProjectsRequest{
		StartEpoch: aYearAgo.Unix(), 
		EndEpoch: now.Unix(),
	}

	err := c.newRequest("POST", "/get-by-date", request, &projects)
	return projects, err
}

// Get tasks for the week by project
func (c *Client) GetWeekTasksByProject(projectId uint) ([]Task, error) {
	var tasks []Task
	
	now := time.Now()

	lastWeek := now.AddDate(0, 0, -7)
	nextWeek := now.AddDate(0, 0, 7)

	request := GetTasksByProjectRequest{
		ProjectId: projectId,
		StartEpoch: lastWeek.Unix(), 
		EndEpoch: nextWeek.Unix(),
	}

	err := c.newRequest("POST", "/tasks/get-by-project", request, &tasks)
	return tasks, err
}

// Add new task
func (c *Client) InsertTask(task Task) error {

	var deadlineEpoch *int64
	var estimDurationMins *int64

	if task.Deadline != nil {
		epoch := task.Deadline.Unix()
		deadlineEpoch = &epoch
	}

	if task.EstimatedDuration != nil {
		mins := int64(task.EstimatedDuration.Minutes())
		estimDurationMins = &mins
	}

	request := InsertTaskRequest{
		ProjectID: task.ProjectID,
		Name: task.Name,
		Priority: task.Priority,
		Description: task.Description,
		Completed: task.Completed,
		DeadlineEpoch: deadlineEpoch,
		EstimatedDurationMins: estimDurationMins,
	}

	var response InsertTaskResponse
	return c.newRequest("POST", "/tasks", request, &response)
}

func (c *Client) UpdateTaskCompletion(id uint, completed bool) error {
	request := UpdateTaskRequest{
		ID: id,
		Completed: completed,
	}

	var response UpdateTaskResponse
	c.newRequest("PATCH", "/tasks/completed", request, &response)

	if response.Success == false {
		return fmt.Errorf(*response.Error)
	}

	return nil
}

// Upsert focus session
func (c *Client) UpsertFocusSession(session FocusSession) (uuid.UUID, error) {
	startEpoch := session.Start.Unix()
	var endEpoch *int64
	
	if(session.End != nil) {
		endTimestamp := session.End.Unix()
		endEpoch = &endTimestamp
	}

	request := UpsertFocusSessionRequest{
		ID: session.ID,
		Device: c.DeviceName,
		CaptureMode: "pomodoro-tui-v2",
		Session: session.Session,
		TaskID: session.TaskID,
		StartEpoch: startEpoch,
		EndEpoch: endEpoch,
	}

	response := new(UpsertFocusSessionResponse)
	err := c.newRequest("POST", "/tasks/focus", request, &response)
	
	sessionId, err := uuid.Parse(response.ID)
	if err != nil {
		log.Fatalf("error parsing API response: %v", err)
	}

	return sessionId, err
}
