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
