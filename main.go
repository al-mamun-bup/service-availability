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

	// Load initial city data (replace with actual city IDs)
	// cityIDs := []int{1, 2, 3, 4, 5, 6, 7} // Add relevant city IDs here

	// for _, cityID := range cityIDs {
	// 	if err := registryService.LoadCitiesFromRegistry(cityID); err != nil {
	// 		log.Printf("Failed to load city data for city ID %d: %v", cityID, err)
	// 	} else {
	// 		log.Printf("Successfully loaded city data for city ID %d", cityID)
	// 	}
	// }

	// Create a new Echo instance
	e := echo.New()

	// Define routes
	e.GET("/v1/is_service_available", availabilityHandler.IsServiceAvailable)

	// Start the server
	log.Println("Server is running on :8080")
	log.Fatal(e.Start(":8080"))
}
