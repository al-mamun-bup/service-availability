package handlers

import (
	"strconv"
)

type DefaultGeoService struct{}

func (d *DefaultGeoService) PointInPolygon(lat, long float64, polygon [][]string) bool {
	var points [][2]float64
	for _, coord := range polygon {
		latVal, _ := strconv.ParseFloat(coord[0], 64)
		longVal, _ := strconv.ParseFloat(coord[1], 64)
		points = append(points, [2]float64{latVal, longVal})
	}

	inside := false
	j := len(points) - 1
	for i := 0; i < len(points); i++ {
		latI, longI := points[i][0], points[i][1]
		latJ, longJ := points[j][0], points[j][1]

		if (longI > long) != (longJ > long) && lat < (latJ-latI)*(long-longI)/(longJ-longI)+latI {
			inside = !inside
		}
		j = i
	}
	return inside
}
