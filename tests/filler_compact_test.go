package tests

import (
	"testing"

	"service-availability/internal/utils"
	h3 "github.com/uber/h3-go/v4"
)

// Helper polygon (square area near Dhaka)
var polygon = []h3.LatLng{
	h3.NewLatLng(23.7808875, 90.2792371),
	h3.NewLatLng(23.7808875, 90.2892371),
	h3.NewLatLng(23.7708875, 90.2892371),
	h3.NewLatLng(23.7708875, 90.2792371),
}

// Test FillPolygonWithHexes returns expected count and values
func TestFillPolygonWithHexes(t *testing.T) {
	resolution := 8

	hexes, err := utils.FillPolygonWithHexes(polygon,resolution)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(hexes) == 0 {
		t.Error("Expected non-empty hex list, got 0")
	}

	for _, hex := range hexes {
		
		if !hex.IsValid() {
			t.Error("Found invalid hex cell")
		}
	}
}

// Test FillPolygonWithHexes with empty polygon
func TestFillPolygonWithHexesEmptyPolygon(t *testing.T) {
	hexes, err := utils.FillPolygonWithHexes([]h3.LatLng{},8)
	if err == nil {
		t.Error("Expected error for empty polygon, got nil")
	}
	if hexes != nil {
		t.Errorf("Expected nil hexes, got: %v", hexes)
	}
}

// Test CompactHexagons compacts a known list correctly
func TestCompactHexagons(t *testing.T) {
	// Generate hexes for a polygon first
	hexes, err := utils.FillPolygonWithHexes(polygon, 8)
	if err != nil {
		t.Fatalf("Error filling hexes: %v", err)
	}

	// Call CustomCompact which now returns an error as well
	compact := utils.CustomCompact(hexes,8)
	// Ensure the compacted result is less than or equal to the original size
	if len(compact) > len(hexes) {
		t.Errorf("Expected compacted size <= original, got %d > %d", len(compact), len(hexes))
	}
}


// Test CompactHexagons with empty input
func TestCompactHexagonsEmpty(t *testing.T) {
	// We expect no error from CustomCompact, so we only check the result.
	hexes := utils.CustomCompact([]h3.Cell{}, 8)

	// Ensure the result is an empty slice
	if len(hexes) != 0 {
		t.Errorf("Expected empty result, got: %v", hexes)
	}
}

// simulate error path if you modify logic to allow injection/mocking
func TestFillPolygonWithHexesInvalidResolution(t *testing.T) {
	_, err := utils.FillPolygonWithHexes(polygon,-1)
	if err == nil {
		t.Error("Expected error for invalid resolution, got nil")
	}
}
