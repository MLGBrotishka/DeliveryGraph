package lstruct

type Coordinate struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Courier struct {
	ID       int        `json:"id"`
	Position Coordinate `json:"position"`
}

type Point struct {
	Lon    float64 `json:"lon"`
	Lat    float64 `json:"lat"`
	Chunks []int
}

func IsCorrectCoordinate(coordinate Coordinate) int {
	if (coordinate.Lon < -180 || coordinate.Lon > 180) {
		return 1
	} else if (coordinate.Lat < -90 || coordinate.Lat > 90) {
		return 2
	}
	return 0
}

func IsCorrectCorier(courier Courier) int {
	if (courier.ID < 0) {
		return 1
	} else if (IsCorrectCoordinate(courier.Position) != 0) {
		return IsCorrectCoordinate(courier.Position)+1
	}
	return 0
}
