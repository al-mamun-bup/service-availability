package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"service-availability/internal/models"
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
		for _, point := range geofence {
			if pair, ok := point.([]interface{}); ok && len(pair) == 2 {
				lat := fmt.Sprintf("%v", pair[0])
				long := fmt.Sprintf("%v", pair[1])
				citySettings.FoodGeofence = append(citySettings.FoodGeofence, []string{lat, long})
			}
		}
	}

	if hours, ok := result["food_open_hours"].([]interface{}); ok {
		for _, h := range hours {
			citySettings.FoodOpenHours = append(citySettings.FoodOpenHours, fmt.Sprintf("%v", h))
		}
	}

	return citySettings, nil
}
