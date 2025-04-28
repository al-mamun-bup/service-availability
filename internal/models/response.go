package models
import h3 "github.com/uber/h3-go/v4"

type CitySettings struct {
	CityID        int        `json:"city_id"`
	FoodGeofence  [][]string `json:"food_geofence"`
	FoodOpenHours []string   `json:"food_open_hours"`
	FoodHexagons  []h3.Cell  `json:"h3_indexes,omitempty"`

}
