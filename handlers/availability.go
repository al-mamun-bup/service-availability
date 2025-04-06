package handlers

import (
	"net/http"
	"service-availability/services"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// type GeoService interface {
// 	PointInPolygon(lat, long float64, polygon [][]string) bool
// }

// type AvailabilityHandler struct {
// 	RegistryService *services.RegistryService
// 	GeoService      GeoService
// }

// func NewAvailabilityHandler(service *services.RegistryService, geoService GeoService) *AvailabilityHandler {
// 	return &AvailabilityHandler{

//			RegistryService: service,
//			GeoService:      geoService,
//		}
//	}
type GeoService interface {
	PointInPolygon(lat, long float64, polygon [][]string) bool
}

type AvailabilityHandler struct {
	RegistryService services.RegistryServiceInterface
	GeoService      GeoService
}

func NewAvailabilityHandler(service services.RegistryServiceInterface, geoService GeoService) *AvailabilityHandler {
	return &AvailabilityHandler{
		RegistryService: service,
		GeoService:      geoService,
	}
}

func (h *AvailabilityHandler) IsServiceAvailable(c echo.Context) error {
	checkType := c.QueryParam("check")
	cityIDStr := c.QueryParam("city_id")
	latStr := c.QueryParam("lat")
	longStr := c.QueryParam("long")

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

	foodOpenHours, err := h.RegistryService.GetFoodOpenHours(cityID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch food open hours"})
	}

	currentTimeUTC := time.Now().UTC()
	offsetSeconds, _ := strconv.Atoi(cityData.TimeOffset)
	offset := time.Duration(offsetSeconds) * time.Second
	localTime := currentTimeUTC.Add(offset)

	currentDay := int(localTime.Weekday())
	openHourData := foodOpenHours[currentDay]
	timeSlots := strings.Split(openHourData, ",")

	if len(timeSlots) < 2 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid food open hours format"})
	}

	openTime, _ := strconv.Atoi(timeSlots[0])
	closeTime, _ := strconv.Atoi(timeSlots[1])
	currentTimeInt := localTime.Hour()*100 + localTime.Minute()
	isTimeAvailable := currentTimeInt >= openTime && currentTimeInt <= closeTime

	geofence, err := services.FetchCityGeofence(cityID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch geofence data"})
	}

	isWithinGeofence := h.GeoService.PointInPolygon(lat, long, geofence.FoodGeofence)
	isAvailable := isTimeAvailable && isWithinGeofence

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Success": isAvailable,
	})
}
