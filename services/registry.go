package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// City structure
type City struct {
	ID            int      `json:"city_id"`
	Name          string   `json:"name"`
	Country       string   `json:"country_name"`
	TimeOffset    int      `json:"time_offset"`
	FoodOpenHours []string `json:"food_open_hours"` // Store open hours
}

// RegistryServiceInterface defines the methods for the registry service
type RegistryServiceInterface interface {
	LoadCitiesFromRegistry() error
	GetCityByID(cityID int) (City, bool)
	GetFoodOpenHours(cityID int) ([]string, error)
}

// RegistryService struct implementing the interface
type RegistryService struct {
	cities      map[int]City
	registryURL string
}

// NewRegistryService initializes the service
func NewRegistryService() *RegistryService {
	service := &RegistryService{
		cities:      make(map[int]City),
		registryURL: getRegistryURL(),
	}
	// Load city data initially
	if err := service.LoadCitiesFromRegistry(); err != nil {
		log.Printf("Error loading cities: %v", err)
	}
	return service
}

// getRegistryURL retrieves the API URL from environment variables
func getRegistryURL() string {
	if url := os.Getenv("REGISTRY_API_URL"); url != "" {
		return url
	}
	return "https://food-registry-v2.p-stageenv.xyz"
}

// LoadCitiesFromRegistry fetches city data from an external API
func (r *RegistryService) LoadCitiesFromRegistry() error {
	url := fmt.Sprintf("%s/api/v1/cities", r.registryURL)
	log.Printf("Fetching city data from: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("request error while fetching cities: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d while fetching cities", resp.StatusCode)
	}

	var result struct {
		Cities []int           `json:"cities"`
		Info   map[string]City `json:"info"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error decoding city data: %w", err)
	}

	for _, cityID := range result.Cities {
		if city, exists := result.Info[fmt.Sprintf("%d", cityID)]; exists {
			r.cities[cityID] = city
		}
	}

	log.Println("City data loaded successfully.")
	return nil
}

// GetFoodOpenHours fetches food open hours dynamically from the API
func (r *RegistryService) GetFoodOpenHours(cityID int) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/settings/cities/%d", r.registryURL, cityID)
	log.Printf("Fetching food open hours from: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request error while fetching food open hours: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d while fetching food open hours", resp.StatusCode)
	}

	var result struct {
		FoodOpenHours []string `json:"food_open_hours"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding food open hours: %w", err)
	}

	return result.FoodOpenHours, nil
}

// GetCityByID retrieves a city by ID
func (r *RegistryService) GetCityByID(cityID int) (City, bool) {
	city, exists := r.cities[cityID]
	return city, exists
}
