package model

type EmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
