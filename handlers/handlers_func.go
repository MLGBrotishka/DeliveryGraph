package handlers

import (
	"encoding/json"
	"fmt"
	"graph/database"
	"graph/lstruct"
	"log"
	"math"
	"net/http"
	"sort"
)

func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatal(err)
	}
}

func AStar(vertices *lstruct.Vertices, edges *lstruct.Edges, startID int, goalID int, max float64, chunks *map[lstruct.Chunk]bool) ([]lstruct.Coordinate, float64) {
	openSet := make(map[int]bool)
	cameFrom := make(map[int]int)
	gScore := make(map[int]float64)
	fScore := make(map[int]float64)

	for id := range *vertices {
		gScore[id] = math.Inf(1)
		fScore[id] = math.Inf(1)
	}
	gScore[startID] = 0
	fScore[startID] = heuristicCost((*vertices)[startID].X, (*vertices)[startID].Y, (*vertices)[goalID].X, (*vertices)[goalID].Y)

	openSet[startID] = true

	for len(openSet) > 0 {
		current := getLowestFScore(&openSet, fScore, gScore, max)
		if current < 0 {
			return nil, math.Inf(1)
		}

		if current == goalID {
			return reconstructPath(vertices, cameFrom, current, startID), gScore[current]
		}

		delete(openSet, current)

		for neighbor := range (*edges)[current] {
			for i := 0; i < len((*vertices)[neighbor].Chunks); i++ {
				_, ok := (*chunks)[(*vertices)[neighbor].Chunks[i]]
				if !ok {
					database.GetVerticesRedis((*vertices)[neighbor].Chunks[i].X, (*vertices)[neighbor].Chunks[i].Y, vertices)
					database.GetEdgesRedis((*vertices)[neighbor].Chunks[i].X, (*vertices)[neighbor].Chunks[i].Y, edges)
					(*chunks)[(*vertices)[neighbor].Chunks[i]] = true
					fmt.Println((*vertices)[neighbor].Chunks[i])
				}
			}

			tentativeGScore := gScore[current] + (*edges)[current][neighbor]
			var ok1 bool
			_, ok1 = gScore[neighbor]
			if !ok1 {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = gScore[neighbor] + heuristicCost((*vertices)[neighbor].X, (*vertices)[neighbor].Y, (*vertices)[goalID].X, (*vertices)[goalID].Y)
				openSet[neighbor] = true
			} else if tentativeGScore < gScore[neighbor] {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = gScore[neighbor] + heuristicCost((*vertices)[neighbor].X, (*vertices)[neighbor].Y, (*vertices)[goalID].X, (*vertices)[goalID].Y)
				openSet[neighbor] = true
			}
		}
	}
	return nil, math.Inf(1)
}

func heuristicCost(aX float64, aY float64, bX float64, bY float64) float64 {
	dx := (aX - bX) * 300
	dy := (aY - bY) * 300
	return dx*dx + dy*dy
}

func getLowestFScore(openSet *map[int]bool, fScore map[int]float64, gScore map[int]float64, max float64) int {
	lowestID := -1
	lowestValue := math.Inf(1)
	for id := range *openSet {
		if max > 0.0 {
			if gScore[id] < max && fScore[id] < lowestValue {
				lowestID = id
				lowestValue = fScore[id]
			} else {
				delete(*openSet, id)
			}
		} else if fScore[id] < lowestValue {
			lowestID = id
			lowestValue = fScore[id]
		}
	}
	return lowestID
}

func reconstructPath(vertices *lstruct.Vertices, cameFrom map[int]int, current int, startID int) []lstruct.Coordinate {
	path := []lstruct.Coordinate{}
	for current != startID {
		vertex := (*vertices)[current]
		path = append([]lstruct.Coordinate{{Lon: vertex.X, Lat: vertex.Y}}, path...)
		current = cameFrom[current]
	}
	vertex := (*vertices)[startID]
	path = append([]lstruct.Coordinate{{Lon: vertex.X, Lat: vertex.Y}}, path...)
	return path
}

func findPoint(x, y float64, vertices *lstruct.Vertices) int {
	var min float64
	var minID int
	minID = -1
	f := 0
	for id := range *vertices {
		p := heuristicCost(x, y, (*vertices)[id].X, (*vertices)[id].Y)
		if f == 0 {
			min = p
			minID = id
			f = 1
		} else if p < min {
			min = p
			minID = id
		}
	}
	return minID
}

func FindChunk(pointX float64, pointY float64) []lstruct.Chunk {
	x := (pointX - database.CenterPoint.Lon) / database.ChunkSize.Lon
	y := (pointY - database.CenterPoint.Lat) / database.ChunkSize.Lat

	if x == float64(int(x)) && y == float64(int(y)) {
		//kek (x-1, y-1), (x, y-1), (x-1, y), (x, y)
		return []lstruct.Chunk{{X: int(x) - 1, Y: int(y) - 1}, {X: int(x), Y: int(y) - 1}, {X: int(x) - 1, Y: int(y)}, {X: int(x), Y: int(y)}}
	} else if x == float64(int(x)) {
		//kek (x-1, int(y)), (x, int(y))
		return []lstruct.Chunk{{X: int(x) - 1, Y: int(y)}, {X: int(x), Y: int(y)}}
	} else if y == float64(int(y)) {
		//kek (int(x), y-1), (int(x), y)
		return []lstruct.Chunk{{X: int(x), Y: int(y) - 1}, {X: int(x), Y: int(y)}}
	} else {
		return []lstruct.Chunk{{X: int(x), Y: int(y)}}
	}
}

func findPath(a lstruct.Coordinate, b lstruct.Coordinate, vertices *lstruct.Vertices, edges *lstruct.Edges, chunks *map[lstruct.Chunk]bool) ([]lstruct.Coordinate, float64) {
	chunksArr := FindChunk(a.Lon, a.Lat)
	for i := 0; i < len(chunksArr); i++ {
		_, ok := (*chunks)[chunksArr[i]]
		if !ok {
			database.GetVerticesRedis(chunksArr[i].X, chunksArr[i].Y, vertices)
			database.GetEdgesRedis(chunksArr[i].X, chunksArr[i].Y, edges)
			(*chunks)[chunksArr[i]] = true
		}
	}
	fmt.Println(chunksArr)
	chunksArr = FindChunk(b.Lon, b.Lat)
	for i := 0; i < len(chunksArr); i++ {
		_, ok := (*chunks)[chunksArr[i]]
		if !ok {
			database.GetVerticesRedis(chunksArr[i].X, chunksArr[i].Y, vertices)
			database.GetEdgesRedis(chunksArr[i].X, chunksArr[i].Y, edges)
			(*chunks)[chunksArr[i]] = true
		}
	}
	fmt.Println(chunksArr)
	startID := findPoint(a.Lon, a.Lat, vertices)
	goalID := findPoint(b.Lon, b.Lat, vertices)
	path, cost := AStar(vertices, edges, startID, goalID, -1.0, chunks)
	return path, cost
}

func sortByHeuristic(points []lstruct.CourierPointID, goal int, vertices *lstruct.Vertices) {
	sort.Slice(points, func(i, j int) bool {
		return heuristicCost((*vertices)[points[i].PointID].X, (*vertices)[points[i].PointID].Y, (*vertices)[goal].X, (*vertices)[goal].Y) < heuristicCost((*vertices)[points[j].PointID].X, (*vertices)[points[j].PointID].Y, (*vertices)[goal].X, (*vertices)[goal].Y)
	})
}

func findClosest(couriers []lstruct.Courier, goal lstruct.Coordinate, vertices *lstruct.Vertices, edges *lstruct.Edges, chunks *map[lstruct.Chunk]bool) ([]lstruct.Coordinate, float64, int) {
	var pointsID []lstruct.CourierPointID
	var chunksArr []lstruct.Chunk
	for i := 0; i < len(couriers); i++ {
		chunksArr = FindChunk(couriers[i].Position.Lon, couriers[i].Position.Lat)
		for i := 0; i < len(chunksArr); i++ {
			_, ok := (*chunks)[chunksArr[i]]
			if !ok {
				database.GetVerticesRedis(chunksArr[i].X, chunksArr[i].Y, vertices)
				database.GetEdgesRedis(chunksArr[i].X, chunksArr[i].Y, edges)
				(*chunks)[chunksArr[i]] = true
			}
		}

		pointsID = append(pointsID, lstruct.CourierPointID{ID: couriers[i].ID, PointID: findPoint(couriers[i].Position.Lon, couriers[i].Position.Lat, vertices)})
	}

	chunksArr = FindChunk(goal.Lon, goal.Lat)
	for i := 0; i < len(chunksArr); i++ {
		_, ok := (*chunks)[chunksArr[i]]
		if !ok {
			database.GetVerticesRedis(chunksArr[i].X, chunksArr[i].Y, vertices)
			database.GetEdgesRedis(chunksArr[i].X, chunksArr[i].Y, edges)
			(*chunks)[chunksArr[i]] = true
		}
	}
	goalID := findPoint(goal.Lon, goal.Lat, vertices)

	sortByHeuristic(pointsID, goalID, vertices)
	var path []lstruct.Coordinate
	cost := -1.0
	id := -1
	ind := -1
	for i := 0; i < len(pointsID); i++ {
		path, cost = AStar(vertices, edges, pointsID[i].PointID, goalID, -1.0, chunks)
		id = pointsID[i].ID
		if path != nil {
			ind = i
			break
		}
	}
	if ind == -1 {
		return nil, math.Inf(1), -1
	} else {
		var path1 []lstruct.Coordinate
		cost1 := -1.0
		id1 := -1
		for i := ind + 1; i < len(pointsID); i++ {
			path1, cost1 = AStar(vertices, edges, pointsID[i].PointID, goalID, cost, chunks)
			id1 = pointsID[i].ID
			if path1 != nil && cost1 < cost {
				path = path1
				cost = cost1
				id = id1
			}
		}
		return path, cost, id
	}
}
