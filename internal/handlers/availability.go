package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	h3 "github.com/uber/h3-go/v4"

	"service-availability/internal/services"
	"service-availability/internal/utils"
)

type ServiceAvailabilityResponse struct {
	Success bool `json:"Success"`
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

	// Time check
	now := time.Now()
	isOpen := utils.IsWithinOpenHours(settings.FoodOpenHours, now)

	isInside := false
	for _, hexagon := range settings.FoodHexagons {
		cell := hexagon
		for res := 7; res <= 9; res++ {
			pointH3,_ := h3.LatLngToCell(h3.LatLng{Lat: lat, Lng: lng}, res)
		// Check if the point (lat, long) is inside the current hexagon
		if pointH3.String() == cell {
			isInside = true
			break
		}
	  }
	  if isInside == true{
		break
	  }
   }
	// Final check
	if isOpen && isInside {
		return c.JSON(http.StatusOK, ServiceAvailabilityResponse{Success: true})
	}

	return c.JSON(http.StatusOK, ServiceAvailabilityResponse{Success: false})
}
