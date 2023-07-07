package lstruct

type Coordinate struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Courier struct {
	ID       int        `json:"id"`
	Position Coordinate `json:"position"`
}

type PathInfo struct {
	CourierID int          `json:"courier-id"`
	Path      []Coordinate `json:"path"`
	Time      int          `json:"time"`
	Cost      float64      `json:"cost"`
}

func IsCorrect(coordinate Coordinate) bool {
	if (coordinate.Lon < -180 || coordinate.Lon > 180 || coordinate.Lat < -90 || coordinate.Lat > 90) {
		return false
	}
	return true
}

func IsCorrect(courier Courier) bool {
	if (courier.ID < 0 || !IsCorrect(courier.Position)) {
		return false
	}
	return true
}
