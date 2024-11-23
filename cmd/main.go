package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"server/internal"
	"server/internal/database"
)

func main() {

	log.Println("Starting application...")

	// PostgreSQL-Verbindung herstellen
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := database.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Tabellen erstellen
	if err := database.Migrate(ctx, dbPool); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Starting webserver...")
	app := fiber.New()

	// Register routes
	internal.RegisterRoutes(app, dbPool)

	// Start server
	err = app.Listen(":80")
	if err != nil {
		return
	}
}
