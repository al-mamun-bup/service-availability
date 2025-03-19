package main

import (
	"service-availability/handlers"
	"service-availability/services"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Initialize the service and handler
	registryService := services.NewRegistryService()
	availabilityHandler := handlers.NewAvailabilityHandler(registryService)

	// Updated Route
	e.GET("/v1/is_service_available", availabilityHandler.IsServiceAvailable)

	// Start the server
	e.Start(":8080")
}
