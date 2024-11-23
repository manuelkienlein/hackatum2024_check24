package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"log"
	"server/internal"
	"server/internal/database"
)

func main() {
	log.Println("Starting application...")

	// Define the dropTables flag
	dropTables := flag.Bool("dropTables", false, "Drop the tables before starting the application")
	flag.Parse()

	// PostgreSQL connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := database.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Drop tables if the flag is set
	if *dropTables {
		if err := database.DropTables(ctx, dbPool); err != nil {
			log.Fatalf("Failed to drop tables: %v", err)
		}
		log.Println("Tables dropped successfully")
	}

	// Migrate the database
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