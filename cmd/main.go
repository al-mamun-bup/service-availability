package main

import (
	"service-availability/internal/handlers"
	"service-availability/config"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadConfig()
	e := echo.New()
	// Route
	e.GET("/v1/is_service_available", handlers.CheckServiceAvailability)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
