package model

type DxoMark struct {
	Name         string `json:"name"`
	Brand        string `json:"brand"`
	Camera       int    `json:"mobileScore"`
	Selfie       int    `json:"selfieScore"`
	Audio        int    `json:"audioScore"`
	Display      int    `json:"displayScore"`
	Battery      int    `json:"batteryScore"`
	Mobile       Mobile `json:"mobile"`
	Price        int    `json:"launch_price"`
	SegmentPrice string `json:"segmentPrice"`
	ImageLink    string `json:"image"`
	LaunchDate   string `json:"launch_date"`
}

type Mobile struct {
	Subscores Subscores `json:"subscores"`
}

type Subscores struct {
	Photo   int `json:"photo"`
	Bokeh   int `json:"bokeh"`
	Preview int `json:"preview"`
	Zoom    int `json:"zoom"`
	Video   int `json:"video"`
}

type DxoMarkRequest struct {
	SortBy string
}
