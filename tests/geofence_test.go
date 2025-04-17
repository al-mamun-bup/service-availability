package tests

import (
	"service-availability/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPointInsidePolygon(t *testing.T) {
	// Define a square
	polygon := [][]string{
		{"0", "0"},
		{"0", "10"},
		{"10", "10"},
		{"10", "0"},
		{"0", "0"}, // closed polygon
	}

	tests := []struct {
		name     string
		point    utils.Point
		expected bool
	}{
		{
			name:     "Point inside square",
			point:    utils.Point{Lat: 5, Lng: 5},
			expected: true,
		},
		{
			name:     "Point outside square",
			point:    utils.Point{Lat: -1, Lng: -1},
			expected: false,
		},
		{
			name:     "Point on edge of square",
			point:    utils.Point{Lat: 0, Lng: 5},
			expected: true, // Ray casting counts edge as inside
		},
		{
			name:     "Malformed coordinate - skip bad point",
			point:    utils.Point{Lat: 5, Lng: 5},
			expected: true, // One malformed won't break others
		},
	}

	// Add a malformed point
	polygonWithBadPoint := append([][]string{{"bad", "coord"}}, polygon...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Malformed coordinate - skip bad point" {
				assert.Equal(t, tt.expected, utils.IsPointInsidePolygon(tt.point, polygonWithBadPoint))
			} else {
				assert.Equal(t, tt.expected, utils.IsPointInsidePolygon(tt.point, polygon))
			}
		})
	}
}
