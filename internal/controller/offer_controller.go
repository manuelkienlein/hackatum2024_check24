package controller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"server/internal/service"
)

type OfferController struct {
	service *service.OfferService
}

// NewOfferController erstellt einen neuen Controller.
func NewOfferController(service *service.OfferService) *OfferController {
	return &OfferController{service: service}
}

// DeleteOffersHandler verarbeitet die DELETE-Anfrage.
func (oc *OfferController) DeleteOffersHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	if err := oc.service.CleanUpOldOffers(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot delete old offers"})
	}

	return c.Status(fiber.StatusOK).SendString("Old offers were cleaned up successfully")
}
