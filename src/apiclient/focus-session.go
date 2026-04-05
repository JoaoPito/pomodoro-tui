package apiclient

import (
	"time"
	"github.com/google/uuid"
)

type FocusSession struct {
	ID		*uuid.UUID	`json:"id"`	
	Start		time.Time	`json:"start_epoch"`
	End		*time.Time	`json:"end_epoch"`
	TaskID 		uint		`json:"task_id"`
	Session 	string		`json:"session"`
}
