package cosine

import "math"

type CosineSimilarityService struct{}

func NewCosineSimilarityService() *CosineSimilarityService {
	return &CosineSimilarityService{}
}

func dotProduct(vec1, vec2 []float64) float64 {
	var dotProduct float64
	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
	}
	return dotProduct
}

// Function to calculate magnitude of a vector
func magnitude(vec []float64) float64 {
	var sumSquares float64
	for _, val := range vec {
		sumSquares += val * val
	}
	return math.Sqrt(sumSquares)
}

// Function to calculate cosine similarity between two vectors
func (css *CosineSimilarityService) CosineSimilarity(vec1, vec2 []float64) float64 {
	return dotProduct(vec1, vec2) / (magnitude(vec1) * magnitude(vec2))
}
