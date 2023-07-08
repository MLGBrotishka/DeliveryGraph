package handlers

import (
	"encoding/json"
	"graph/lstruct"
	"net/http"
)

func GetV1PathMultipleCouriers(w http.ResponseWriter, r *http.Request) {
	var pathMSRequest lstruct.PathMultipleStartRequest
	err := json.NewDecoder(r.Body).Decode(&pathMSRequest)
	if err != nil {
		errorResponse := lstruct.ErrorResponse{
			Message: "Bad Input",
		}
		SendJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	response := lstruct.ErrorResponse{
		Message: "Not implemented",
	}

	SendJSONResponse(w, http.StatusOK, response)
}
