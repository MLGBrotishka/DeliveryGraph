package lstruct

type PathRequest struct {
	Courier       Courier    `json:"courier"`
	EndCoordinate Coordinate `json:"end_coordinate"`
}

type PathMultipleStartRequest struct {
	Couriers      []Courier  `json:"couriers"`
	EndCoordinate Coordinate `json:"end_coordinate"`
}
