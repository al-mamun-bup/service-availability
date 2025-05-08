package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"io"
	"net/http"

	"service-availability/internal/models"
	"service-availability/internal/utils"
	h3 "github.com/uber/h3-go/v4"
)

var FetchCitySettings = fetchCitySettings
var RegistryBaseURL = "https://food-registry-v2.p-stageenv.xyz"

func fetchCitySettings(cityID int) (*models.CitySettings, error) {
	url := fmt.Sprintf("%s/api/v1/settings/cities/%d", RegistryBaseURL, cityID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch city settings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response from registry API: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Extract required fields manually
	citySettings := &models.CitySettings{}

	if cityIDVal, ok := result["city_id"].(int); ok {
		citySettings.CityID = int(cityIDVal)
	}

	if geofence, ok := result["food_geofence"].([]interface{}); ok {
		var boundary []h3.LatLng
		for _, point := range geofence {
			if pair, ok := point.([]interface{}); ok && len(pair) == 2 {
			// Handle lat and long correctly with type assertion and conversion
			var lat, lng float64
			latStr := fmt.Sprintf("%v", pair[0])
			lngStr := fmt.Sprintf("%v", pair[1])

			lat, err = strconv.ParseFloat(latStr, 64)
			if err != nil {
    			return nil, fmt.Errorf("invalid lat value: %v", pair[0])
			}

			lng, err = strconv.ParseFloat(lngStr, 64)
			if err != nil {
    			return nil, fmt.Errorf("invalid long value: %v", pair[1])
			}

				// Add lat, long to boundary for H3 LatLng
				boundary = append(boundary, h3.NewLatLng(lat, lng))
				// Store lat, long as strings in FoodGeofence for later use (if needed)
				citySettings.FoodGeofence = append(citySettings.FoodGeofence, []string{fmt.Sprintf("%v", lat), fmt.Sprintf("%v", lng)})
			}
		}
	
		// Convert polygon to hexagons
		hexagons, err := utils.FillPolygonWithHexes(boundary)
		if err != nil {
			return nil, fmt.Errorf("failed to generate hexagons: %w", err)
		}
	
		// Store hexagons in city settings
		citySettings.FoodHexagons = hexagons
	}
	

	if hours, ok := result["food_open_hours"].([]interface{}); ok {
		for _, h := range hours {
			citySettings.FoodOpenHours = append(citySettings.FoodOpenHours, fmt.Sprintf("%v", h))
		}
	}

	return citySettings, nil
}
