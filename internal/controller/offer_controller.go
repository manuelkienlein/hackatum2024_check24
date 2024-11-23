package controller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"server/internal/models"
	"server/internal/service"
)

type OfferController struct {
	offerService *service.OfferService
}

// NewOfferController erstellt einen neuen Controller.
func NewOfferController(service *service.OfferService) *OfferController {
	return &OfferController{offerService: service}
}

// CreateOffersHandler verarbeitet die POST-Anfragen
func (oc *OfferController) CreateOffersHandler(c *fiber.Ctx) error {
	var request struct {
		Offers []models.Offer `json:"offers"`
	}

	// Parse the request body
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Unable to parse JSON: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// Call service to create offers
	if err := oc.offerService.CreateOffers(c.Context(), request.Offers); err != nil {
		log.Printf("Error creating offers: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create offers"})
	}

	return c.Status(fiber.StatusOK).SendString("Offers were created successfully")
}

// DeleteOffersHandler verarbeitet die DELETE-Anfrage.
func (oc *OfferController) DeleteOffersHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	if err := oc.offerService.CleanUpOldOffers(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot delete old offers"})
	}

	return c.Status(fiber.StatusOK).SendString("Old offers were cleaned up successfully")
}
