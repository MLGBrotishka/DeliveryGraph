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
	if request.Message != "iamgay" && request.Message != "iamgayestgay" {
		errorResponse := lstruct.ErrorResponse{
			Message: "Wrong message",
		}
		SendJSONResponse(w, http.StatusOK, errorResponse)
		return
	}
	// Создание и отправка ответа
	if request.Message == "iamgayestgay" {
		database.EraseAllTablesRedis()
		response := lstruct.ErrorResponse{
			Message: "All tables cleared successfully",
		}
		SendJSONResponse(w, http.StatusOK, response)
		return
	}
	if request.Message == "iamgay" {
		database.LoadFromTextToRedis("./database/offline/chunks")
		database.FindCollisions()
		response := lstruct.ErrorResponse{
			Message: "Example loaded successfully",
		}
		SendJSONResponse(w, http.StatusOK, response)
		return
	}

}
