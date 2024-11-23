package main

import (
	"github.com/gofiber/fiber/v2"
	"server/internal"
)

func main() {
	// Initialize
	dbPool := internal.InitDB()
	app := fiber.New()

	// Register routes
	internal.RegisterRoutes(app, dbPool)

	// Start server
	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
