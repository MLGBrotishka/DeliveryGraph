package lstruct

type Coordinate struct {
	Lon float64 `json:"lon" validate:"required"`
	Lat float64 `json:"lat" validate:"required"`
}

type Courier struct {
	ID       int        `json:"id" validate:"required"`
	Position Coordinate `json:"position" validate:"required"`
}

type CourierPointID struct {
	ID      int
	PointID int
}

type Vertices map[int]Vertex

type Vertex struct {
	X      float64
	Y      float64
	Chunks []Chunk
}

type Chunk struct {
	X int
	Y int
}

type Edges map[int]map[int]float64

func IsCorrectCoordinate(coordinate Coordinate) int {
	if coordinate.Lon < -180 || coordinate.Lon > 180 {
		return 1
	} else if coordinate.Lat < -90 || coordinate.Lat > 90 {

		return 2
	}
	return 0
}

func IsCorrectCourier(courier Courier) int {
	if courier.ID < 0 {
		return 1
	} else if IsCorrectCoordinate(courier.Position) != 0 {
		return IsCorrectCoordinate(courier.Position) + 1
	}
	return 0
}
