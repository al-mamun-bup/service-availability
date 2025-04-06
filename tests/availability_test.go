package tests

import (
	"net/http"
	"net/http/httptest"
	"service-availability/handlers"
	"service-availability/services"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRegistryService mocks the RegistryService for testing
type MockRegistryService struct {
	mock.Mock
}

func (m *MockRegistryService) GetCityByID(cityID int) (services.City, bool) {
	args := m.Called(cityID)
	return args.Get(0).(services.City), args.Bool(1)
}

func (m *MockRegistryService) LoadCitiesFromRegistry(cityID int) error {
	args := m.Called(cityID)
	return args.Error(0)
}

func (m *MockRegistryService) GetFoodOpenHours(cityID int) ([]string, error) {
	args := m.Called(cityID)
	return args.Get(0).([]string), args.Error(1)
}

// MockGeoService mocks the GeoService for testing
type MockGeoService struct {
	mock.Mock
}

func (m *MockGeoService) PointInPolygon(lat, long float64, polygon [][]string) bool {
	args := m.Called(lat, long, polygon)
	return args.Bool(0)
}

func TestIsServiceAvailable(t *testing.T) {
	// Initialize mock services
	mockRegistryService := new(MockRegistryService) // Mock of RegistryService interface
	mockGeoService := new(MockGeoService)

	// Create an instance of AvailabilityHandler with mock services
	availabilityHandler := handlers.NewAvailabilityHandler(mockRegistryService, mockGeoService)

	// Mock data for the city
	cityID := 1
	city := services.City{
		ID:         cityID,
		Name:       "New York",
		Country:    "USA",
		TimeOffset: "3600", // 1 hour offset
	}
	foodOpenHours := []string{"1000,2200", "0900,2100", "1100,2300", "1000,2200", "0900,2100", "1100,2300", "1000,2200"}

	// Setup mock behaviors for the registry service
	mockRegistryService.On("GetCityByID", cityID).Return(city, true)
	mockRegistryService.On("GetFoodOpenHours", cityID).Return(foodOpenHours, nil)
	mockRegistryService.On("LoadCitiesFromRegistry", cityID).Return(nil)

	// Setup mock behavior for the geo service
	mockGeoService.On("PointInPolygon", 40.7128, -74.0060, mock.Anything).Return(true)

	// Create a mock request and response recorder
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/is_service_available?check=food&city_id=1&lat=40.7128&long=-74.0060", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler function
	if assert.NoError(t, availabilityHandler.IsServiceAvailable(c)) {
		// Assert response status
		assert.Equal(t, http.StatusOK, rec.Code)

		// Assert the response body
		assert.Contains(t, rec.Body.String(), `"Success": true`)
	}

	// Assert that the mock methods were called
	mockRegistryService.AssertExpectations(t)
	mockGeoService.AssertExpectations(t)
}
