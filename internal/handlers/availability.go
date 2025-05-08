package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	h3 "github.com/uber/h3-go/v4"

	"service-availability/internal/services"
	"service-availability/internal/utils"
	"service-availability/config"
)

type ServiceAvailabilityResponse struct {
	Success bool `json:"Success"`
}


func IsPointInsideHexagons(lat, lng float64, hexagons []string, minResolution, maxResolution int) bool {
	for _, hex := range hexagons {
		for res := minResolution; res <= maxResolution; res++ {
			pointH3, err := h3.LatLngToCell(h3.LatLng{Lat: lat, Lng: lng}, res)
			if err != nil {
				continue
			}
			if pointH3.String() == hex {
				return true
			}
		}
	}
	return false
}

func CheckServiceAvailability(c echo.Context) error {
	check := c.QueryParam("check")
	cityIDStr := c.QueryParam("city_id")
	latStr := c.QueryParam("lat")
	longStr := c.QueryParam("lng")

	if check != "food" || cityIDStr == "" || latStr == "" || longStr == "" {
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

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lng, err2 := strconv.ParseFloat(longStr, 64)
	if err1 != nil || err2 != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid lat or long format",
		})
	}

	settings, err := services.FetchCitySettings(cityID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch city settings",
		})
	}

	// Check time-based availability
	now := time.Now()
	isOpen := utils.IsWithinOpenHours(settings.FoodOpenHours, now)

	// Check geofence-based availability
	min := config.AppConfig.H3.ResolutionMin
    max := config.AppConfig.H3.ResolutionMax
	
	isInside := IsPointInsideHexagons(lat, lng, settings.FoodHexagons, min, max)

	// Final response
	return c.JSON(http.StatusOK, ServiceAvailabilityResponse{
		Success: isOpen && isInside,
	})
}
