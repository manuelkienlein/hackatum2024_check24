package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

func RegisterRoutes(app *fiber.App, dbPool *pgxpool.Pool) {

	app.Get("/api/offers", func(c *fiber.Ctx) error {
		return FilterOffers(dbPool, c)
	})

	app.Post("/api/offers", func(c *fiber.Ctx) error {
		return CreateOffers(dbPool, c)
	})

}
