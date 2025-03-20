package handlers

import (
	"net/http"
	"service-availability/services"
	"strconv"
	"strings"
	"time"

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

	// Get city details
	cityData, exists := h.RegistryService.GetCityByID(cityID)
	if !exists {
		if err := h.RegistryService.LoadCitiesFromRegistry(); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to refresh city data"})
		}
		cityData, exists = h.RegistryService.GetCityByID(cityID)
		if !exists {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "City not found"})
		}
	}

	// Fetch food open hours for the city
	foodOpenHours, err := h.RegistryService.GetFoodOpenHours(cityID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch food open hours"})
	}

	// Get current UTC time and apply offset
	currentTimeUTC := time.Now().UTC()
	offsetHours := cityData.TimeOffset / 3600
	offsetMinutes := (cityData.TimeOffset % 3600) / 60
	localTime := currentTimeUTC.Add(time.Duration(offsetHours)*time.Hour + time.Duration(offsetMinutes)*time.Minute)

	// Get current day index (Sunday = 0, Monday = 1, ..., Saturday = 6)
	currentDay := int(localTime.Weekday())

	// Parse open hours for the current day
	openHourData := foodOpenHours[currentDay]
	timeSlots := strings.Split(openHourData, ",")

	if len(timeSlots) < 2 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid food open hours format"})
	}

	// Convert open-close time range from string to int
	openTime, _ := strconv.Atoi(timeSlots[0])  // Opening time in HHMM format
	closeTime, _ := strconv.Atoi(timeSlots[1]) // Closing time in HHMM format

	// Extract current hour and minute
	currentTimeInt := localTime.Hour()*100 + localTime.Minute()

	// Check if current time is within the open hours
	isAvailable := currentTimeInt >= openTime && currentTimeInt <= closeTime

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Success": isAvailable,
		//"LocalTime": localTime.Format("2006-01-02 15:04:05"),
	})
}
