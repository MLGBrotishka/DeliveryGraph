package handlers

import (
	"encoding/json"
	"fmt"
	"graph/database"
	"graph/lstruct"
	"log"
	"net/http"
)

func GetV1Path(w http.ResponseWriter, r *http.Request) {
	var courier lstruct.Courier
	err := json.NewDecoder(r.Body).Decode(&courier)
	if err != nil {
		errorResponse := lstruct.ErrorResponse{
			Message: "Bad Input",
		}
		SendJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	res := lstruct.IsCorrectCorier(courier)
	if res != 0 {
		var errorResponse lstruct.ErrorResponse
		if res == 1 {
			errorResponse = lstruct.ErrorResponse{
				Message: "Uncorrect courier ID",
			}
		} else if res == 2 {
			errorResponse = lstruct.ErrorResponse{
				Message: "Longitude out of range",
			}
		} else if res == 3 {
			errorResponse = lstruct.ErrorResponse{
				Message: "Latitude out of range",
			}
		}
		SendJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}
	database.SelectRedis(courier.ID)
	log.Println(database.GetRedis(fmt.Sprintf("%v", courier.Position.Lat)))

	// Создание и отправка ответа
	response := lstruct.ErrorResponse{
		Message: "Not implemented",
	}

	SendJSONResponse(w, http.StatusOK, response)
}
