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

	// Validate input
	if checkType != "food" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid check type"})
	}

	cityID, err := strconv.Atoi(cityIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid city_id"})
	}

	// Check if city exists
	cityData, exists := h.RegistryService.GetCityByID(cityID)
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "City not found"})
	}

	// Here, we simply return service availability as "available" for demo purposes
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "Service is available",
		"check_type": checkType,
		"city_id":    cityID,
		"city_name":  cityData.Name,
		"timezone":   cityData.Timezone,
	})
}
