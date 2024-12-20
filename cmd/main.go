package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"server/internal/controller"
	"server/internal/database"
	"server/internal/framework"
	"server/internal/repository"
	"server/internal/service"
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

	// Init components
	offerRepo := repository.NewOfferRepository(dbPool)
	offerService := service.NewOfferService(offerRepo)
	offerController := controller.NewOfferController(offerService)

	log.Println("Starting webserver...")
	app := fiber.New()

	// Add logger
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${url} | ${error}\n",
	}))

	// Register new routes
	framework.RegisterRoutes(app, offerController)

	// Add swagger
	framework.RegisterSwagger(app)

	// Start server
	err = app.Listen(":80")
	if err != nil {
		return
	}
}
