package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type City struct {
	ID            int        `json:"city_id"`
	Name          string     `json:"name"`
	Country       string     `json:"country_name"`
	TimeOffset    string     `json:"time_offset"`
	FoodOpenHours []string   `json:"food_open_hours"`
	FoodGeofence  [][]string `json:"food_geofence"`
}

type RegistryService struct {
	cities map[int]City
}

func NewRegistryService() *RegistryService {
	return &RegistryService{
		cities: make(map[int]City),
	}
}

func (r *RegistryService) LoadCitiesFromRegistry(cityID int) error {
	url := fmt.Sprintf("https://food-registry-v2.p-stageenv.xyz/api/v1/settings/cities/%d", cityID)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("failed to fetch city data from registry API")
	}
	defer resp.Body.Close()

	var result struct {
		CityID int `json:"city_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error decoding city data: %v", err)
	}

	// Add the city ID to the internal map
	r.cities[result.CityID] = City{
		ID: result.CityID,
		// You can populate other fields of City here as necessary.
	}

	fmt.Println("City data loaded successfully.")
	return nil
}

// Fetch food_open_hours and food_geofence dynamically from the API
func (r *RegistryService) GetFoodOpenHours(cityID int) ([]string, error) {
	url := fmt.Sprintf("https://food-registry-v2.p-stageenv.xyz/api/v1/settings/cities/%d", cityID)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch food open hours from registry API")
	}
	defer resp.Body.Close()

	var result struct {
		FoodOpenHours []string   `json:"food_open_hours"`
		FoodGeofence  [][]string `json:"food_geofence"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding food open hours: %v", err)
	}

	// Update geofence data for the city
	if city, exists := r.cities[cityID]; exists {
		city.FoodGeofence = result.FoodGeofence
		r.cities[cityID] = city
	}

	return result.FoodOpenHours, nil
}

func (r *RegistryService) GetCityByID(cityID int) (City, bool) {
	city, exists := r.cities[cityID]
	return city, exists
}

// FetchCityGeofence returns the geofence for a given city
func FetchCityGeofence(cityID int) (City, error) {
	url := fmt.Sprintf("https://food-registry-v2.p-stageenv.xyz/api/v1/settings/cities/%d", cityID)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {

		return City{}, errors.New("failed to fetch city geofence from registry API")
	}

	defer resp.Body.Close()

	var result City
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println(err)
		return City{}, fmt.Errorf("error decoding city geofence: %v", err)
	}

	return result, nil
}
