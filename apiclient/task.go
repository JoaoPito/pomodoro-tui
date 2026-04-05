package apiclient

import (
	"time"
)

type Task struct {
	ID			uint		`json:"id"`
	ProjectID		uint		`json:"project_id"`
	Priority		uint		`json:"priority"`
	Name			string		`json:"name"`
	Completed		bool		`json:"completed"`
	Description		string		`json:"description"`
	Deadline		*time.Time	`json:"deadline"`
	EstimatedDuration	*time.Duration	`json:"estimated_duration_min"`
}
