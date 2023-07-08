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
	
	vertices := Vertices{ }
	edges := Edges{ }
	chunks := map[Chunk]bool { }

	path, cost, id := findClosest(pathMSRequest.Couriers, pathMSRequest.EndCoordinate, &vertices, &edges, &chunks)
	if path != nil {
		response := PathInfoResponse{
			id,
			path,
			cost,
			cost
		}
	} else {
		response := lstruct.ErrorResponse{
			Message: "Kuda blyat",
		}

	}

	SendJSONResponse(w, http.StatusOK, response)
}
