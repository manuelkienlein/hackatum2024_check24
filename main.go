package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/api/offers", func(c *fiber.Ctx) error {
		return c.SendString("GET Offers!")
	})

	app.Post("/api/offers", func(c *fiber.Ctx) error {
		// endpoint is expected to be blocking. Your solution should only return a response once the offers have been successfully processed and can be read from with new GET requests. During this blocking time, there are no requests to the GET endpoint we don't send any request that would require the offer to be present or change the expected result.
		return c.SendString("POST Offers!")
	})

	app.Delete("/api/offers", func(c *fiber.Ctx) error {
		return c.SendString("DELETE Offers!")
	})

	app.Listen(":3000")
}
