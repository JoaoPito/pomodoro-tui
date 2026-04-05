package apiclient

import "time"

type Project struct {
	ID 		uint		`json:"id"`
	Name 		string		`json:"name"`
	Created_at 	time.Time	`json:"created_at"`
	Repository 	string		`json:"repository"`
	Creator 	string		`json:"creator"`
	Archived 	bool		`json:"archived"`
	Updated_at 	time.Time	`json:"updated_at"`
	Last_task_updated_at 	time.Time	`json:"last_task_updated_at"`
}
