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
	latStr := c.QueryParam("lat")
	longStr := c.QueryParam("long")

	// Validate required parameters
	if checkType != "food" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid check type"})
	}

	cityID, err := strconv.Atoi(cityIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid city_id"})
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid latitude"})
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid longitude"})
	}

	// Get city details
	cityData, exists := h.RegistryService.GetCityByID(cityID)
	if !exists {
		if err := h.RegistryService.LoadCitiesFromRegistry(cityID); err != nil {
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
	x, err := strconv.Atoi(cityData.TimeOffset)
	y, err := strconv.Atoi(cityData.TimeOffset)
	offsetHours := x / 3600
	offsetMinutes := (y % 3600) / 60
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
	isTimeAvailable := currentTimeInt >= openTime && currentTimeInt <= closeTime

	// Fetch geofence data from the registry service
	geofence, err := services.FetchCityGeofence(cityID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch geofence data"})
	}

	// Check if the point is within the geofence
	isWithinGeofence := PointInPolygon(lat, long, geofence.FoodGeofence)

	// Both conditions must be true
	isAvailable := isTimeAvailable && isWithinGeofence

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Success": isAvailable,
	})
}

// PointInPolygon checks if a point is inside a polygon using the ray-casting algorithm
func PointInPolygon(lat, long float64, polygon [][]string) bool {
	var points [][2]float64
	for _, coord := range polygon {
		latVal, _ := strconv.ParseFloat(coord[0], 64)
		longVal, _ := strconv.ParseFloat(coord[1], 64)
		points = append(points, [2]float64{latVal, longVal})
	}

	inside := false
	j := len(points) - 1
	for i := 0; i < len(points); i++ {
		latI, longI := points[i][0], points[i][1]
		latJ, longJ := points[j][0], points[j][1]

		if (longI > long) != (longJ > long) && lat < (latJ-latI)*(long-longI)/(longJ-longI)+latI {
			inside = !inside
		}
		j = i
	}
	return inside
}
