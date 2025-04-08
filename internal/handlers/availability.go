package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"service-availability/internal/services"
	"service-availability/internal/utils"
)

type ServiceAvailabilityResponse struct {
	Success bool `json:"Success"`
}

func CheckServiceAvailability(c echo.Context) error {
	check := c.QueryParam("check")
	cityIDStr := c.QueryParam("city_id")

	if check != "food" || cityIDStr == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Missing or invalid query parameters",
		})
	}

	cityID, err := strconv.Atoi(cityIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid city_id format",
		})
	}

	// Fetch city settings
	settings, err := services.FetchCitySettings(cityID)
	if err != nil {
		c.Logger().Errorf("fetch error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch city settings",
		})
	}

	// Time check
	now := time.Now()
	isOpen := utils.IsWithinOpenHours(settings.FoodOpenHours, now)

	if isOpen {
		return c.JSON(http.StatusOK, ServiceAvailabilityResponse{Success: true})
	}

	return c.JSON(http.StatusOK, ServiceAvailabilityResponse{Success: false})
}
