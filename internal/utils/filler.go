package utils

import (
	"fmt"
	h3 "github.com/uber/h3-go/v4"
)

// FillPolygonWithHexes fills a polygon with uncompacted hexagons at a fixed resolution.
func FillPolygonWithHexes(boundary []h3.LatLng) ([]string, error) {
	hexes := PolygonToH3Indexes(boundary)
	fmt.Print(hexes)
	return hexes, nil
}
