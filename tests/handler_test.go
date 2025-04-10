package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"service-availability/internal/handlers"
	"service-availability/internal/models"
	"service-availability/internal/services"
)

// Inject mock into services package
func init() {
	services.FetchCitySettings = func(cityID int) (*models.CitySettings, error) {
		mockCitySettings := models.CitySettings{
			CityID: 1,
			FoodOpenHours: []string{
				"0800,2359,08:00AM-23:59PM",
			},
			FoodGeofence: [][]string{
				{"22.37928754838105", "91.8551245973755"},
				{"22.382648561771163", "91.85044232911378"},
				{"22.3831921", "91.8450305"},
				{"22.379801", "91.844310"},
				{"22.377000", "91.847000"},
				{"22.376200", "91.850800"},
				{"22.37928754838105", "91.8551245973755"},
			},
		}
		if cityID == 1 {
			return &mockCitySettings, nil
		}
		return &models.CitySettings{}, echo.NewHTTPError(http.StatusNotFound, "City not found")
	}
}

func TestCheckServiceAvailability(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success: inside geofence and open hours",
			queryParams:    "?check=food&city_id=1&lat=22.3795&long=91.8510",
			expectedStatus: http.StatusOK,
			expectedBody:   `"Success":true`,
		},
		{
			name:           "Failure: outside geofence but within open hours",
			queryParams:    "?check=food&city_id=1&lat=22.3650&long=91.8600",
			expectedStatus: http.StatusOK,
			expectedBody:   `"Success":false`,
		},
		{
			name:           "Failure: invalid city_id",
			queryParams:    "?check=food&city_id=abc&lat=22.3795&long=91.8510",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"Invalid city_id format"`,
		},
		{
			name:           "Failure: missing parameters",
			queryParams:    "?check=food&city_id=1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"Missing or invalid query parameters"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/is_service_available"+tc.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// ✅ Call standalone function directly
			err := handlers.CheckServiceAvailability(c)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.expectedStatus, rec.Code)
				assert.Contains(t, rec.Body.String(), tc.expectedBody)
			}
		})
	}
}
