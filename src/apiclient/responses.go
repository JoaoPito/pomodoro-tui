package apiclient

type UpsertFocusSessionResponse struct {
	ID string `json:"id"`
}

type InsertTaskResponse struct {
	ID uint `json:"id"`
}

type UpdateTaskResponse struct {
	Success	bool	`json:"success"`
	ID 	uint 	`json:"id"`
	Error	*string	`json:"error"`
}

type DeleteTaskResponse struct {
	Success bool    `json:"success"`
	ID      uint    `json:"id"`
	Error   *string `json:"error"`
}
