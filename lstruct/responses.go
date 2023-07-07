package lstruct

type PathInfo struct {
	CourierID int          `json:"courier-id"`
	Path      []Coordinate `json:"path"`
	Time      int          `json:"time"`
	Cost      float64      `json:"cost"`
}

type PointAvailableResponse struct {
	Available bool `json:"available"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
