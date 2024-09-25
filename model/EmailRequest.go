package model

type EmailRequest struct {
	UserEmail       string           `json:"email"`
	Recommendations []Recommendation `json:"recommendations"`
}

type Recommendation struct {
	Name  string `json:"name"`
	Liked bool   `json:"liked"`
}
