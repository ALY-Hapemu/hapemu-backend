package service

import (
	"encoding/json"
	"hapemu/accessor"
	"hapemu/cosine"
	"hapemu/model"
	"hapemu/vector"
	"net/http"
	"sort"

	_ "github.com/lib/pq"
)

type RecommendationService struct {
	cosineSimilarityService *cosine.CosineSimilarityService
	vectorGeneratorService  *vector.VectorGeneratorService
	hapemuDatabaseAccessor  *accessor.HapemuDatabaseAccessor
}

func NewRecommendationService(cs *cosine.CosineSimilarityService, vgs *vector.VectorGeneratorService, hda *accessor.HapemuDatabaseAccessor) *RecommendationService {
	return &RecommendationService{
		cosineSimilarityService: cs,
		vectorGeneratorService:  vgs,
		hapemuDatabaseAccessor:  hda,
	}
}

// main function
func (rs *RecommendationService) RecommendSmartphones(w http.ResponseWriter, r *http.Request) {
	var recommendationsRequest model.RecommendationsRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&recommendationsRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userPreferenceVector = rs.vectorGeneratorService.ConvertUserPreferenceToTargetVec(recommendationsRequest)
	var similarities []model.SmartphoneSimilarity
	var smartphones = rs.hapemuDatabaseAccessor.GetSmartphoneList()

	for _, smartphone := range smartphones {
		var smartphonesVec = rs.vectorGeneratorService.ConvertSmartphoneToVec(smartphone, userPreferenceVector)
		similarity := rs.cosineSimilarityService.CosineSimilarity(smartphonesVec, userPreferenceVector)
		similarities = append(similarities, model.SmartphoneSimilarity{Name: smartphone.Name, Similarity: similarity})
	}
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	if len(similarities) > 5 {
		similarities = similarities[:5]
	}
	var recommendationsResponse model.RecommendationsResponse
	for _, similarity := range similarities {
		recommendationsResponse.Recommendations = append(recommendationsResponse.Recommendations, similarity.Name)
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(recommendationsResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
