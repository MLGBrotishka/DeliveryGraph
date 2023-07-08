package handlers

import (
	"encoding/json"
	"graph/lstruct"
	"net/http"
)

func GetV1Path(w http.ResponseWriter, r *http.Request) {
	var pathRequest lstruct.PathRequest
	err := json.NewDecoder(r.Body).Decode(&pathRequest)
	if err != nil {
		errorResponse := lstruct.ErrorResponse{
			Message: "Bad Input",
		}
		SendJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	res := lstruct.IsCorrectCourier(pathRequest.Courier)
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

	vertices := lstruct.Vertices{}
	edges := lstruct.Edges{}
	chunks := map[lstruct.Chunk]bool{}

	path, cost := findPath(pathRequest.Courier.Position, pathRequest.EndCoordinate, &vertices, &edges, &chunks)

	// Создание и отправка ответа
	if path != nil {
		response := lstruct.PathInfoResponse{
			CourierID: pathRequest.Courier.ID,
			Path:      path,
			Time:      int(cost),
			Cost:      cost,
		}
		SendJSONResponse(w, http.StatusOK, response)
	} else {
		response := lstruct.ErrorResponse{
			Message: "Not reachable from point destination",
		}
		SendJSONResponse(w, http.StatusNotAcceptable, response)
	}

}
