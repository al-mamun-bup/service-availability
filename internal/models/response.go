package models

type CitySettings struct {
	CityID        int        `json:"city_id"`
	FoodGeofence  [][]string `json:"food_geofence"`
	FoodOpenHours []string   `json:"food_open_hours"`
}
