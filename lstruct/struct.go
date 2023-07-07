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
