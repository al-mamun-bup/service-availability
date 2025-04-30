package utils

import (
	"fmt"
	h3 "github.com/uber/h3-go/v4"
)

// FillPolygonWithHexes fills a polygon with uncompacted hexagons at a fixed resolution.
func FillPolygonWithHexes(boundary []h3.LatLng) ([]string, error) {
	// Create a GeoPolygon from boundary
	// polygon := h3.GeoPolygon{
	// 	GeoLoop: boundary,
	// 	Holes:   nil,
	// }
	// geoPoly := h3.GeoPolygon{
    //     GeoLoop: boundary,
    // }

	// Fill polygon with compacted hexes
	hexes := PolygonToH3Indexes(boundary)
	//compacted:= customPolyFillUsingRange(geoPoly, 7,10)
	fmt.Print(hexes)
	return hexes, nil
}
