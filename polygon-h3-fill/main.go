package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"polygon-h3-fill/polygon"
	"polygon-h3-fill/utils"
	h3 "github.com/uber/h3-go/v4"
)

type GeoJSONPolygon struct {
	Coordinates [][]float64 `json:"coordinates"`
}

func main() {
	// Load polygon from JSON
	file, err := os.ReadFile("data/polygon.json")
	if err != nil {
		log.Fatalf("failed to read polygon file: %v", err)
	}

	var geo GeoJSONPolygon
	if err := json.Unmarshal(file, &geo); err != nil {
		log.Fatalf("failed to parse polygon JSON: %v", err)
	}

	// Convert coordinates to []h3.LatLng
	var boundary []h3.LatLng
	for _, coord := range geo.Coordinates {
		latlng := h3.NewLatLng(coord[1], coord[0])
		boundary = append(boundary, latlng)
	}

	// Set the resolution
	res := 10

	// Fill polygon with H3 hexagons
	hexes, err := polygon.FillPolygonWithHexes(boundary, res)
	if err != nil {
		log.Fatalf("failed to fill polygon: %v", err)
	}

	// fmt.Println(hexes)

	// Compact the hexes
	//compacted, err := h3.CompactCells(hexes)
	compacted := utils.CustomCompact(hexes, res)

	if err != nil {
		log.Fatalf("failed to compact hexes: %v", err)
	}

	// Print compacted hexes
	for _, h := range compacted {
		hx := uint64(h)
		fmt.Println(h3.IndexToString(hx))
	}
	fmt.Println("Total compacted hexes:", len(compacted))
}
