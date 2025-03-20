package main

import (
	"log"
	"service-availability/handlers"
	"service-availability/services"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize the registry service and the handler
	registryService := services.NewRegistryService()
	availabilityHandler := handlers.NewAvailabilityHandler(registryService)

	// Create a new Echo instance
	e := echo.New()

	// Define routes
	e.GET("/v1/is_service_available", availabilityHandler.IsServiceAvailable)

	// Start the server
	log.Fatal(e.Start(":8080"))
}
