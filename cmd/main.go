package main

import (
	"github.com/gofiber/fiber/v2"
	"server/internal"
)

func main() {
	app := fiber.New()

	internal.RegisterRoutes(app)

	app.Listen(":80")
}
