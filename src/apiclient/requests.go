package apiclient

import "github.com/google/uuid"

type GetLatestProjectsRequest struct {
	StartEpoch int64 `json:"start_epoch"`
	EndEpoch int64 `json:"end_epoch"`
}

type GetTasksByProjectRequest struct {
	StartEpoch int64 `json:"start_epoch"`
	EndEpoch int64 `json:"end_epoch"`
	ProjectId uint `json:"project_id"`
}

type UpsertFocusSessionRequest struct {
	ID		*uuid.UUID	`json:"id,omitempty"`	
	StartEpoch	int64		`json:"start_epoch"`
	EndEpoch	*int64		`json:"end_epoch,omitempty"`
	TaskID 		uint		`json:"task_id"`
	Session 	string		`json:"session"`
	Device	 	string		`json:"device"`
	CaptureMode 	string		`json:"capture_mode"`
}

type InsertTaskRequest struct {
	ProjectID		uint		`json:"project_id"`
	Priority		uint		`json:"priority"`
	Name			string		`json:"name"`
	Completed		bool		`json:"completed"`
	Description		string		`json:"description"`
	DeadlineEpoch		*int64		`json:"deadline_epoch"`
	EstimatedDurationMins	*int64		`json:"estimated_duration_min"`
}

type UpdateTaskRequest struct {
	ID		uint	`json:"id"`
	Completed	bool	`json:"completed"`
}
