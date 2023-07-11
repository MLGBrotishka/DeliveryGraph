package handlers

import (
	"encoding/json"
	"graph/database"
	"graph/lstruct"
	"net/http"
)

func GetV1SecretLoadDatabase(w http.ResponseWriter, r *http.Request) {
	var request lstruct.ErrorResponse
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorResponse := lstruct.ErrorResponse{
			Message: "Bad Input",
		}
		SendJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}
	if request.Message != "clear" && request.Message != "load" && request.Message != "load" {
		errorResponse := lstruct.ErrorResponse{
			Message: "Wrong message",
		}
		SendJSONResponse(w, http.StatusOK, errorResponse)
		return
	}
	// Создание и отправка ответа
	if request.Message == "clear" {
		database.EraseAllTablesRedis()
		response := lstruct.ErrorResponse{
			Message: "All tables cleared successfully",
		}
		SendJSONResponse(w, http.StatusOK, response)
		return
	}
	if request.Message == "load" {
		database.LoadFromTextToRedis("./database/offline/chunks")
		//online.FindCollisions(1, 1)
		response := lstruct.ErrorResponse{
			Message: "Example loaded successfully",
		}
		SendJSONResponse(w, http.StatusOK, response)
		return
	}

}
