package service

import (
	"encoding/json"
	"hapemu/model"
	"net/http"
)

func GetAntutuList(w http.ResponseWriter, r *http.Request) {
	// call antutu python scrape
	// convert json from txt file to antutu golang class
	// return list of antutu

	decoder := json.NewDecoder(r.Body)
	var request model.Smartphone
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
}
