package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hapemu/model"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host     = "aws-0-ap-southeast-1.pooler.supabase.com"
	port     = 5432
	user     = "postgres.yovcppevikilglvpktzq"
	password = "HapemuPostgres123"
	dbname   = "postgres"
)

// region convert smartphone from database to vector
func getValueForPrice(price string) float64 {
	if strings.Compare(price, "essential") == 0 {
		return 1
	} else if strings.Compare(price, "mid") == 0 {
		return 2
	} else if strings.Compare(price, "high") == 0 {
		return 3
	}
	return 4
}

func getValueForCamera(camera int) float64 {
	if camera <= 75 {
		return 1
	} else if camera <= 104 {
		return 2
	} else if camera <= 134 {
		return 3
	}
	return 4
}

func getValueForBattery(battery string) float64 {
	var batteryValue, err = strconv.Atoi(battery)
	if err != nil {
		fmt.Println("error converting battery string to integer")
	}
	if batteryValue < 4000 {
		return 1
	} else if batteryValue < 4500 {
		return 2
	} else if batteryValue < 5000 {
		return 3
	}
	return 4
}

func getVecValueFromRam(ram string) float64 {
	if ram == "1" || ram == "2" || ram == "4" {
		return 1
	} else if ram == "6" || ram == "8" {
		return 2
	} else if ram == "12" {
		return 3
	} else {
		return 4
	}
}

func getValueForRam(ram string, ramVec float64) float64 {
	var ramList = []string{"1", "2", "4", "6", "8", "12", "16", "32"}
	var minVec, maxVec float64
	for _, cur := range ramList {
		if strings.Contains(ram, cur) {
			minVec = getVecValueFromRam(cur)
			break
		}
	}
	for _, cur := range ramList {
		if strings.Contains(ram, cur) {
			maxVec = getVecValueFromRam(cur)
		}
	}
	if ramVec < minVec {
		return minVec
	} else if ramVec >= minVec && ramVec <= maxVec {
		return ramVec
	} else {
		return maxVec
	}
}

func getVecValueFromStorage(storage string) float64 {
	if storage == "32GB" || storage == "64GB" || storage == "128GB" {
		return 1
	} else if storage == "256GB" {
		return 2
	} else if storage == "512GB" {
		return 3
	} else {
		return 4
	}
}

func getValueForStorage(storage string, storageVec float64) float64 {
	var storageList = []string{"32GB", "64GB", "128GB", "256GB", "512GB", "1TB"}
	var minVec, maxVec float64
	for _, cur := range storageList {
		if strings.Contains(storage, cur) {
			minVec = getVecValueFromStorage(cur)
			break
		}
	}
	for _, cur := range storageList {
		if strings.Contains(storage, cur) {
			maxVec = getVecValueFromStorage(cur)
		}
	}
	if storageVec < minVec {
		return minVec
	} else if storageVec >= minVec && storageVec <= maxVec {
		return storageVec
	} else {
		return maxVec
	}
}

func convertSmartphoneToVec(smartphone model.Smartphone, targetVec []float64) []float64 {
	var smartphonesVecs []float64
	smartphonesVecs = append(smartphonesVecs, getValueForPrice(smartphone.SegmentPrice)) // price
	// smartphonesVecs = append(smartphonesVecs) // processor
	smartphonesVecs = append(smartphonesVecs, getValueForCamera(smartphone.DxomarkScore))           // camera
	smartphonesVecs = append(smartphonesVecs, getValueForBattery(smartphone.Battery))               // battery
	smartphonesVecs = append(smartphonesVecs, getValueForRam(smartphone.Ram, targetVec[4]))         // ram
	smartphonesVecs = append(smartphonesVecs, getValueForStorage(smartphone.Storage, targetVec[5])) // storage
	return smartphonesVecs
}

// endregion

// region apply cosine similarity algorithm
// Function to calculate dot product of two vectors
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
func cosineSimilarity(vec1, vec2 []float64) float64 {
	return dotProduct(vec1, vec2) / (magnitude(vec1) * magnitude(vec2))
}

// Function to recommend movies based on cosine similarity
func recommendSmartphone(smartphones []model.Smartphone, targetPhoneVec []float64) []model.SmartphoneSimilarity {
	var similarities []model.SmartphoneSimilarity

	// Calculate similarity of each movie with the target movie
	for _, smartphone := range smartphones {
		similarity := cosineSimilarity(convertSmartphoneToVec(smartphone, targetPhoneVec), targetPhoneVec)
		similarities = append(similarities, model.SmartphoneSimilarity{Name: smartphone.Name, Similarity: similarity})
	}

	// Sort movies by similarity in descending order
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	// Return top 5 similar movies
	if len(similarities) > 5 {
		return similarities[:5]
	}
	return similarities
}

// endregion

// region convert user quiz to vector
func getPriceValue(price string) float64 {
	if strings.Contains(price, "essensial") {
		return 1
	} else if strings.Contains(price, "midrange") {
		return 2
	} else if strings.Contains(price, "premium") {
		return 3
	}
	return 4
}

func getValue(str string) float64 {
	if strings.Contains(str, "tidak") {
		return 1
	} else if strings.Contains(str, "cukup") {
		return 2
	} else if strings.Contains(str, "penting") {
		return 3
	}
	return 4
}

func convertRecommendationRequestToTargetVec(request model.RecommendationsRequest) []float64 {
	var vec []float64

	vec = append(vec, getPriceValue(request.Price))
	vec = append(vec, getValue(request.Processor))
	vec = append(vec, getValue(request.Camera))
	vec = append(vec, getValue(request.Baterry))
	vec = append(vec, getValue(request.Ram))
	vec = append(vec, getValue(request.Storage))

	return vec
}

//endregion

// main function
func RecommendSmartphones(w http.ResponseWriter, r *http.Request) {
	var recommendationsRequest model.RecommendationsRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&recommendationsRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var targetPhoneVec = convertRecommendationRequestToTargetVec(recommendationsRequest)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var smartphones []model.Smartphone
	sqlStatement := `SELECT name, segmentPrice, processor, dxomarkScore, battery, ram, storage FROM "Smartphone"`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var smartphone model.Smartphone
		err := rows.Scan(&smartphone.Name, &smartphone.SegmentPrice, &smartphone.Processor, &smartphone.DxomarkScore, &smartphone.Battery, &smartphone.Ram, &smartphone.Storage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		smartphones = append(smartphones, smartphone)
	}

	var similarities = recommendSmartphone(smartphones, targetPhoneVec)
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