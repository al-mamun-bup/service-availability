package main

import (
	"log"
	"service-availability/handlers"
	"service-availability/services"

	"github.com/labstack/echo/v4"
)

func main() {
	registryService := services.NewRegistryService()
	geoService := &handlers.DefaultGeoService{}
	availabilityHandler := handlers.NewAvailabilityHandler(registryService, geoService)

	e := echo.New()
	e.GET("/v1/is_service_available", availabilityHandler.IsServiceAvailable)

	log.Println("Server is running on :8080")
	log.Fatal(e.Start(":8080"))
}
