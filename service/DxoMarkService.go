package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hapemu/model"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "aws-0-ap-southeast-1.pooler.supabase.com"
	port     = 5432
	user     = "postgres.yovcppevikilglvpktzq"
	password = "HapemuPostgres123"
	dbname   = "postgres"
	format   = "2006-01-02" // specify yyyy-mm-dd format
)

func InsertDxoMarkToDatabase() {
	// Command to execute Python script with arguments
	cmd := exec.Command("python3", "/Users/t-albertus.dharmala/skripsi/hapemu-scrape/dxomark.py") // change to the one in your local

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running Python script: %v", err)
	}

	// Path to the output file
	outputFilePath := "/Users/t-albertus.dharmala/skripsi/hapemu-backend/dxomark.json"
	fmt.Println("Output JSON file saved at:", outputFilePath)

	// Read the JSON file
	jsonFile, err := os.Open(outputFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer jsonFile.Close()

	// Decode the JSON file into a struct
	var dxomarkList []model.DxoMark
	err = json.NewDecoder(jsonFile).Decode(&dxomarkList)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// validate to only insert with camera score
	var filteredDxomarkList []model.DxoMark
	for _, element := range dxomarkList {
		if element.Camera != 0 {
			filteredDxomarkList = append(filteredDxomarkList, element)
		}
	}
	fmt.Printf("total: %d\n", len(filteredDxomarkList))

	// create query for insert list of dxomark
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	sqlDeleteStatement := `DELETE FROM "Smartphone"`
	_, err = db.Exec(sqlDeleteStatement)
	if err != nil {
		log.Fatal("Error Deleting from database:", err)
	}
	sqlStatement := `
	INSERT INTO "Smartphone" (name, brand, "dxomarkScore", photo, bokeh, preview, zoom, video, price, "segmentPrice", "imageLink", "launchDate", "isLastThreeYear")
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING id`

	var smartphoneList []model.Smartphone
	for _, dxomark := range filteredDxomarkList {
		// ToDo: add scraping method from gsm
		// get the result and mapToSmartphone (just add new parameter for scrape result of gsm)
		var smartphone = mapToSmartphone(dxomark)
		smartphoneList = append(smartphoneList, smartphone)
	}

	for _, smartphone := range smartphoneList {
		_, err = db.Exec(sqlStatement, smartphone.Name, smartphone.Brand, smartphone.DxomarkScore, smartphone.Photo, smartphone.Bokeh, smartphone.Preview,
			smartphone.Zoom, smartphone.Video, smartphone.Price, smartphone.SegmentPrice, smartphone.ImageLink, smartphone.LaunchDate, smartphone.IsLastThreeYears)
		if err != nil {
			log.Fatal("Error inserting into database:", err)
		}
	}
}

// to do
func GetSmartphoneList(w http.ResponseWriter, r *http.Request) {
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	log.Fatal("Error connecting to the database:", err)
	// }
	// defer db.Close()
	// sqlStatement := `SELECT * FROM "Smartphone"`
	// rows, err := db.Query(sqlStatement)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// if err != nil {
	// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 	return
	// }

	// var response []model.Smartphone
	// for rows.Next() {
	// 	var smartphone model.Smartphone
	// 	err := rows.Scan(&smartphone.Name, &smartphone.Brand, &smartphone.DxomarkScore, &smartphone.Photo, &smartphone.Bokeh,
	// 		&smartphone.Preview, &smartphone.Zoom, &smartphone.Video, &smartphone.Price, &smartphone.SegmentPrice, &smartphone.ImageLink, &smartphone.LaunchDate)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	response = append(response, smartphone)
	// }

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
}

func mapToSmartphone(dxomark model.DxoMark) model.Smartphone {
	var smartphone model.Smartphone
	smartphone.Name = dxomark.Name
	smartphone.Brand = dxomark.Brand
	smartphone.DxomarkScore = dxomark.Camera
	smartphone.Photo = dxomark.Mobile.Subscores.Photo
	smartphone.Zoom = dxomark.Mobile.Subscores.Zoom
	smartphone.Bokeh = dxomark.Mobile.Subscores.Bokeh
	smartphone.Zoom = dxomark.Mobile.Subscores.Zoom
	smartphone.Video = dxomark.Mobile.Subscores.Video
	smartphone.ImageLink = dxomark.ImageLink
	smartphone.Price = dxomark.Price
	smartphone.SegmentPrice = dxomark.SegmentPrice
	smartphone.LaunchDate = dxomark.LaunchDate
	smartphone.IsLastThreeYears = IsLastThreeYears(smartphone.LaunchDate)
	return smartphone
}

func IsLastThreeYears(dateString string) bool {
	parsedDate, err := time.Parse(format, dateString)
	if err != nil {
		fmt.Println("error when parsing date")
	}
	currentTime := time.Now()
	threeYearsAgo := currentTime.AddDate(-3, 0, 0)

	return parsedDate.After(threeYearsAgo) && parsedDate.Before(currentTime)
}
