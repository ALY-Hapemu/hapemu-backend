package model

type EmailRequest struct {
	UserEmail       string   `json:"email"`
	Recommendations []string `json:"recommendations"`
}
