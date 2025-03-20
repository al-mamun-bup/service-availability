package handlers

import (
	"net/http"
	"service-availability/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AvailabilityHandler struct {
	RegistryService *services.RegistryService
}

func NewAvailabilityHandler(service *services.RegistryService) *AvailabilityHandler {
	return &AvailabilityHandler{
		RegistryService: service,
	}
}

func (h *AvailabilityHandler) IsServiceAvailable(c echo.Context) error {
	checkType := c.QueryParam("check")
	cityIDStr := c.QueryParam("city_id")

	if checkType != "food" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid check type"})
	}

	cityID, err := strconv.Atoi(cityIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid city_id"})
	}

	// Check if the city exists in the cache
	cityData, exists := h.RegistryService.GetCityByID(cityID)
	if !exists {
		// Try to refresh the registry data
		if err := h.RegistryService.LoadCitiesFromRegistry(); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to refresh city data"})
		}

		// Re-check after refreshing
		cityData, exists = h.RegistryService.GetCityByID(cityID)
		if !exists {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "City not found"})
		}
	}

	// Check time_offset based on country_name
	var expectedOffset int
	switch cityData.Country {
	case "Bangladesh":
		expectedOffset = 21600
	case "Nepal":
		expectedOffset = 20700
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unsupported country"})
	}

	// Validate the time_offset
	if cityData.TimeOffset == expectedOffset {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"Success": true,
		})
	}

	// If the time offset doesn't match, service is unavailable
	return c.JSON(http.StatusOK, map[string]interface{}{
		"Success": false,
	})
}
