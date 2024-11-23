package framework

import (
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"server/internal/controller"
)

func RegisterRoutes(app *fiber.App, offerController *controller.OfferController) {
	app.Delete("/api/offers", offerController.DeleteOffersHandler)
	app.Post("/api/offers", offerController.CreateOffersHandler)
	app.Get("/api/offers", offerController.GetOffersHandler)
}

func RegisterSwagger(app *fiber.App) {
	cfg := swagger.Config{
		BasePath: "/",
		FilePath: "./docs/openapi.yaml",
		Path:     "openapi",
		Title:    "Swagger API Docs",
	}

	app.Use(swagger.New(cfg))
}
