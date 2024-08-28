package main

import (
	"fmt"
	"hapemu/accessor"
	"hapemu/cosine"
	"hapemu/service"
	"hapemu/vector"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	cosineSimilarityService := cosine.NewCosineSimilarityService()
	vectorGeneratorService := vector.NewVectorGeneratorService()
	hapemuDatabaseAccessor := accessor.NewHapemuDatabaseAccessor()
	emailService := service.NewEmailService(587, "smtp.gmail.com", "hapemu.id@gmail.com", "fdvtmvobhemhxvvi")
	recommendationService := service.NewRecommendationService(cosineSimilarityService, vectorGeneratorService, hapemuDatabaseAccessor)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	http.HandleFunc("/send-email", emailService.EmailRecommendations)
	http.HandleFunc("/get-recommendations", recommendationService.RecommendSmartphones)

	handler := cors.Default().Handler(http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
