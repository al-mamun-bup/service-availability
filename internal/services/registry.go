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
				var lat, long float64
				switch v := pair[0].(type) {
				case float64:
					lat = v
				case string:
					latConv, err := strconv.ParseFloat(v, 64)
					if err != nil {
						return nil, fmt.Errorf("invalid lat value: %v", v)
					}
					lat = latConv
				default:
					return nil, fmt.Errorf("invalid lat type: %v", v)
				}
	
				switch v := pair[1].(type) {
				case float64:
					long = v
				case string:
					longConv, err := strconv.ParseFloat(v, 64)
					if err != nil {
						return nil, fmt.Errorf("invalid long value: %v", v)
					}
					long = longConv
				default:
					return nil, fmt.Errorf("invalid long type: %v", v)
				}
	
				// Add lat, long to boundary for H3 LatLng
				boundary = append(boundary, h3.NewLatLng(lat, long))
	
				// Store lat, long as strings in FoodGeofence for later use (if needed)
				citySettings.FoodGeofence = append(citySettings.FoodGeofence, []string{fmt.Sprintf("%v", lat), fmt.Sprintf("%v", long)})
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
