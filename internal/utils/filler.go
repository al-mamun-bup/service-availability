package utils

import (
	"fmt"
	h3 "github.com/uber/h3-go/v4"
)

// FillPolygonWithHexes fills a polygon with uncompacted hexagons at a fixed resolution.
func FillPolygonWithHexes(boundary []h3.LatLng, resolution int) ([]h3.Cell, error) {
	// Create a GeoPolygon from boundary
	polygon := h3.GeoPolygon{
		GeoLoop: boundary,
		Holes:   nil,
	}

	// Fill polygon with compacted hexes
	hexes, err := h3.PolygonToCells(polygon, resolution)
	if err != nil {
		return nil, fmt.Errorf("error generating cells: %w", err)
	}

	// Uncompact /Compact (if needed) — idempotent in this case but explicit
	compacted:= CustomCompact(hexes, resolution)

	return compacted, nil
}
