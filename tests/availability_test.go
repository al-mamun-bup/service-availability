package tests

import (
	"net/http"
	"net/http/httptest"
	"service-availability/handlers"
	"service-availability/services"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

// ✅ Mock struct implementing RegistryServiceInterface
type MockRegistryService struct {
	mock.Mock
}

// Implement LoadCitiesFromRegistry for Mock
func (m *MockRegistryService) LoadCitiesFromRegistry() error {
	args := m.Called()
	return args.Error(0)
}

// Implement GetCityByID for Mock
func (m *MockRegistryService) GetCityByID(cityID int) (services.City, bool) {
	args := m.Called(cityID)
	city, ok := args.Get(0).(services.City)
	return city, ok
}

// ✅ Implement GetFoodOpenHours for Mock
func (m *MockRegistryService) GetFoodOpenHours(cityID int) ([]string, error) {
	args := m.Called(cityID)
	return args.Get(0).([]string), args.Error(1)
}

// ✅ Ensure the mock implements the interface
var _ services.RegistryServiceInterface = (*MockRegistryService)(nil)

func TestIsServiceAvailable(t *testing.T) {
	e := echo.New()

	mockService := new(MockRegistryService)
	handler := handlers.NewAvailabilityHandler(mockService) // ✅ Now it works!

	req := httptest.NewRequest(http.MethodGet, "/v1/is_service_available?city_id=1&check=food", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock the expected city data
	mockService.On("GetCityByID", 1).Return(services.City{ID: 1, Country: "Bangladesh", TimeOffset: 21600}, true)
	mockService.On("GetFoodOpenHours", 1).Return([]string{"0900,2200", "1000,2100", "1100,2000", "1200,2300", "0800,1800", "0700,1700", "0600,1600"}, nil)

	// Call the handler
	handler.IsServiceAvailable(c)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", rec.Code)
	}
}
