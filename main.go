package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Coordinate struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Courier struct {
	ID       int        `json:"id"`
	Position Coordinate `json:"position"`
}

type PathInfo struct {
	CourierID int          `json:"courier-id"`
	Path      []Coordinate `json:"path"`
	Time      int          `json:"time"`
	Cost      float64      `json:"cost"`
}

type PointAvailableResponse struct {
	Available bool `json:"available"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

//ручки

func getFastestPathHandler(w http.ResponseWriter, r *http.Request) {
	// Обработка запроса /v1/path
	// Реализуйте здесь логику для получения самого быстрого пути
}

func getFastestPathForMultipleCouriersHandler(w http.ResponseWriter, r *http.Request) {
	// Обработка запроса /v1/path/multiple-couriers
	// Реализуйте здесь логику для получения самого быстрого пути для нескольких курьеров
}

func isPointAvailableHandler(w http.ResponseWriter, r *http.Request) {
	// Чтение и декодирование тела запроса в формате JSON
	var coordinate Coordinate
	err := json.NewDecoder(r.Body).Decode(&coordinate)
	if err != nil {
		errorResponse := ErrorResponse{
			Message: "Bad Input",
		}
		sendJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	// Проверка принадлежности координат диапазону
	available := coordinate.Lon >= -90 && coordinate.Lon <= 90 && coordinate.Lat >= -90 && coordinate.Lat <= 90

	// Создание и отправка ответа
	response := PointAvailableResponse{
		Available: available,
	}

	sendJSONResponse(w, http.StatusOK, response)
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	http.HandleFunc("/v1/path", getFastestPathHandler)
	http.HandleFunc("/v1/path/multiple-couriers", getFastestPathForMultipleCouriersHandler)
	http.HandleFunc("/v1/point/is-available", isPointAvailableHandler)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
