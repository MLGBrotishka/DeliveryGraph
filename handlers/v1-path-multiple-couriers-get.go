package handlers

import (
	"graph/lstruct"
	"net/http"
)

func GetV1PathMultipleCouriers(w http.ResponseWriter, r *http.Request) {
	response := lstruct.ErrorResponse{
		Message: "Not implemented",
	}

	SendJSONResponse(w, http.StatusOK, response)
}
