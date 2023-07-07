package lstruct

type PointAvailableResponse struct {
	Available bool `json:"available"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
