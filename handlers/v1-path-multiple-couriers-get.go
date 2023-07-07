package handlers

import (
	"encoding/json"
	"graph/lstruct"
	"net/http"
)

func GetV1PathMultipleCouriers(w http.ResponseWriter, r *http.Request) {
	//для теста
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
	response := lstruct.ErrorResponse{
		Message: "Not implemented",
	}

	SendJSONResponse(w, http.StatusOK, response)
}
