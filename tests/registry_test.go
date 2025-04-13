package tests

import (
	"errors"
	"service-availability/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCityByID(t *testing.T) {
	// Create an instance of MockRegistryService
	mockRegistryService := new(MockRegistryService)

	// Mock data for the city
	cityID := 1
	city := services.City{
		ID:         cityID,
		Name:       "New York",
		Country:    "USA",
		TimeOffset: "3600",
	}

	// Setup mock behavior for GetCityByID
	mockRegistryService.On("GetCityByID", cityID).Return(city, true)

	// Call GetCityByID method
	result, exists := mockRegistryService.GetCityByID(cityID)

	// Assertions
	assert.True(t, exists)
	assert.Equal(t, cityID, result.ID)
	assert.Equal(t, "New York", result.Name)

	// Assert that the mock method was called
	mockRegistryService.AssertExpectations(t)
}

func TestLoadCitiesFromRegistry(t *testing.T) {
	// Create an instance of MockRegistryService
	mockRegistryService := new(MockRegistryService)

	// Mock the behavior of LoadCitiesFromRegistry to simulate success
	cityID := 1
	mockRegistryService.On("LoadCitiesFromRegistry", cityID).Return(nil)

	// Call LoadCitiesFromRegistry method
	err := mockRegistryService.LoadCitiesFromRegistry(cityID)

	// Assertions
	assert.NoError(t, err)

	// Assert that the mock method was called
	mockRegistryService.AssertExpectations(t)
}

func TestLoadCitiesFromRegistry_Error(t *testing.T) {
	// Create an instance of MockRegistryService
	mockRegistryService := new(MockRegistryService)

	// Mock the behavior of LoadCitiesFromRegistry to simulate an error
	cityID := 1
	mockRegistryService.On("LoadCitiesFromRegistry", cityID).Return(errors.New("failed to fetch city data"))

	// Call LoadCitiesFromRegistry method
	err := mockRegistryService.LoadCitiesFromRegistry(cityID)

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to fetch city data")

	// Assert that the mock method was called
	mockRegistryService.AssertExpectations(t)
}

func TestGetFoodOpenHours(t *testing.T) {
	// Create an instance of MockRegistryService
	mockRegistryService := new(MockRegistryService)

	// Mock data for food open hours
	cityID := 1
	foodOpenHours := []string{"1000,2200", "0900,2100", "1100,2300", "1000,2200", "0900,2100", "1100,2300", "1000,2200"}

	// Setup mock behavior for GetFoodOpenHours
	mockRegistryService.On("GetFoodOpenHours", cityID).Return(foodOpenHours, nil)

	// Call GetFoodOpenHours method
	result, err := mockRegistryService.GetFoodOpenHours(cityID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, foodOpenHours, result)

	// Assert that the mock method was called
	mockRegistryService.AssertExpectations(t)
}

func TestGetFoodOpenHours_Error(t *testing.T) {
	// Create an instance of MockRegistryService
	mockRegistryService := new(MockRegistryService)

	// Setup mock behavior for GetFoodOpenHours to simulate an error
	cityID := 1
	mockRegistryService.On("GetFoodOpenHours", cityID).Return(nil, errors.New("failed to fetch food open hours"))

	// Call GetFoodOpenHours method
	result, err := mockRegistryService.GetFoodOpenHours(cityID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "failed to fetch food open hours")

	// Assert that the mock method was called
	mockRegistryService.AssertExpectations(t)
}
