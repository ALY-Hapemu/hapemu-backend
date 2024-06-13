package main

import (
	"hapemu/service"
)

func main() {
	service.InsertDxoMarkToDatabase()

	// http.HandleFunc("/antutu", service.GetAntutuList)
	// http.HandleFunc("/dxomark", service.GetDxoMarkList)

	// log.Fatal(http.ListenAndServe(":8080", nil))
}
