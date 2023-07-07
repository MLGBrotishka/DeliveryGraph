package main

import (
	"graph/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/v1/path", handlers.GetV1Path)
	http.HandleFunc("/v1/path/multiple-couriers", handlers.GetV1PathMultipleCouriers)
	http.HandleFunc("/v1/point/is-available", handlers.GetV1PointIsAvailable)
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
