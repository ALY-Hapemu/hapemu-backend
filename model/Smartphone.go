package model

// Struct to hold an item and its similarity score
type SmartphoneSimilarity struct {
	Name       string
	Similarity float64
}

type Smartphone struct {
	Name         string
	SegmentPrice string
	Processor    string
	DxomarkScore int
	Battery      string
	Ram          string
	Storage      string
}
