package handlers

import (
	"encoding/json"
	"graph/lstruct"
	"net/http"
	"time"
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

	err = lstruct.ValidatePathMultiple(pathMSRequest)
	if err != nil {
		errorResponse := lstruct.ErrorResponse{
			Message: err.Error(),
		}
		SendJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	vertices := lstruct.Vertices{}
	edges := lstruct.Edges{}
	chunks := map[lstruct.Chunk]bool{}

	path, cost, id := findClosest(pathMSRequest.Couriers, pathMSRequest.EndCoordinate, &vertices, &edges, &chunks)
	rub := float64(cost * GetTimeValue(time.Now().String()) * OkladPerHour(cost) / 3600)
	if path != nil {
		response := lstruct.PathInfoResponse{
			CourierID: id,
			Path:      path,
			Time:      int(cost),
			Cost:      rub,
		}
		SendJSONResponse(w, http.StatusOK, response)
	} else {
		response := lstruct.ErrorResponse{
			Message: "Not reachable from point destination",
		}
		SendJSONResponse(w, http.StatusNotAcceptable, response)
	}
}
