package utils

import (
	"strconv"
)

// Point represents a latitude/longitude point
type Point struct {
	Lat float64
	Lng float64
}

// IsPointInsidePolygon checks if a point is inside a polygon using ray casting
func IsPointInsidePolygon(point Point, polygon [][]string) bool {
	var polyPoints []Point
	for _, coord := range polygon {
		lat, err1 := strconv.ParseFloat(coord[0], 64)
		lng, err2 := strconv.ParseFloat(coord[1], 64)
		if err1 != nil || err2 != nil {
			continue
		}
		polyPoints = append(polyPoints, Point{Lat: lat, Lng: lng})
	}

	n := len(polyPoints)
	inside := false

	j := n - 1
	for i := 0; i < n; i++ {
		if (polyPoints[i].Lng > point.Lng) != (polyPoints[j].Lng > point.Lng) &&
			(point.Lat < (polyPoints[j].Lat-polyPoints[i].Lat)*(point.Lng-polyPoints[i].Lng)/(polyPoints[j].Lng-polyPoints[i].Lng)+polyPoints[i].Lat) {
			inside = !inside
		}
		j = i
	}

	return inside
}
