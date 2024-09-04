package model

type RecommendationsRequest struct {
	Price     string `json:"price"`
	Processor string `json:"processor"`
	Camera    string `json:"camera"`
	Battery   string `json:"battery"`
	Ram       string `json:"ram"`
	Storage   string `json:"storage"`
}
