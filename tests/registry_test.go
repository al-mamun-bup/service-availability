package tests_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"service-availability/internal/models"
	"service-availability/internal/services"

	"github.com/stretchr/testify/assert"
)

func TestFetchCitySettings_Success(t *testing.T) {
	// Mock response data
	mockResponse := models.CitySettings{
		CityID: 1,
		FoodOpenHours: []string{
			"0800,2359,08:00AM-23:59PM",
		},
		FoodGeofence: [][]string{
			{"22.37928754838105", "91.8551245973755"},
			{"22.382648561771163", "91.85044232911378"},
		},
	}

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/api/v1/settings/cities/1")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Override the URL in the actual service for test purposes
	originalBaseURL := services.RegistryBaseURL
	services.RegistryBaseURL = server.URL
	defer func() { services.RegistryBaseURL = originalBaseURL }()

	// Run the test
	resp, err := services.FetchCitySettings(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.CityID)
	assert.Equal(t, "0800,2359,08:00AM-23:59PM", resp.FoodOpenHours[0])
}

func TestFetchCitySettings_NotFound(t *testing.T) {
	// Mock a 404 response from the registry API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "City not found", http.StatusNotFound)
	}))
	defer server.Close()

	// Override the URL in the service for testing
	originalBaseURL := services.RegistryBaseURL
	services.RegistryBaseURL = server.URL
	defer func() { services.RegistryBaseURL = originalBaseURL }()

	resp, err := services.FetchCitySettings(99)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "404")
}
