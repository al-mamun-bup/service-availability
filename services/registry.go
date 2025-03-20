package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type City struct {
	ID            int      `json:"city_id"`
	Name          string   `json:"name"`
	Country       string   `json:"country_name"`
	TimeOffset    int      `json:"time_offset"`
	FoodOpenHours []string `json:"food_open_hours"` // Store open hours
}

type RegistryService struct {
	cities map[int]City
}

func NewRegistryService() *RegistryService {
	service := &RegistryService{
		cities: make(map[int]City),
	}
	// Load city data initially
	service.LoadCitiesFromRegistry()
	return service
}

func (r *RegistryService) LoadCitiesFromRegistry() error {
	url := "https://food-registry-v2.p-stageenv.xyz/api/v1/cities"
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("failed to fetch cities from registry API")
	}
	defer resp.Body.Close()

	var result struct {
		Cities []int           `json:"cities"`
		Info   map[string]City `json:"info"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error decoding city data: %v", err)
	}

	// Update the internal map
	for _, cityID := range result.Cities {
		if city, exists := result.Info[fmt.Sprintf("%d", cityID)]; exists {
			r.cities[cityID] = city
		}
	}

	fmt.Println("City data loaded successfully.")
	return nil
}

// Fetch food_open_hours dynamically from the API
func (r *RegistryService) GetFoodOpenHours(cityID int) ([]string, error) {
	url := fmt.Sprintf("https://food-registry-v2.p-stageenv.xyz/api/v1/settings/cities/%d", cityID)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch food open hours from registry API")
	}
	defer resp.Body.Close()

	var result struct {
		FoodOpenHours []string `json:"food_open_hours"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding food open hours: %v", err)
	}

	return result.FoodOpenHours, nil
}

func (r *RegistryService) GetCityByID(cityID int) (City, bool) {
	city, exists := r.cities[cityID]
	return city, exists
}
