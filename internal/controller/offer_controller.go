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

// GetOffersHandler verarbeitet GET-Abfragen
/*func (oc *OfferController) GetOffersHandler(c *fiber.Ctx) error {
	// TODO

	// Parse query parameters into the struct
	params := models.OfferFilterParams{
		RegionID:              c.QueryInt("regionID"),
		TimeRangeStart:        c.QueryInt("timeRangeStart"),
		TimeRangeEnd:          c.QueryInt("timeRangeEnd"),
		NumberDays:            c.QueryInt("numberDays"),
		SortOrder:             c.Query("sortOrder"),
		Page:                  c.QueryInt("page"),
		PageSize:              c.QueryInt("pageSize"),
		PriceRangeWidth:       c.QueryInt("priceRangeWidth"),
		MinFreeKilometerWidth: c.QueryInt("minFreeKilometerWidth"),
		MinNumberSeats:        c.QueryInt("minNumberSeats", 0),
		MinPrice:              c.QueryInt("minPrice", 0),
		MaxPrice:              c.QueryInt("maxPrice", 0),
		CarType:               c.Query("carType"),
		OnlyVollkasko:         c.Query("onlyVollkasko") == "true",
		MinFreeKilometer:      c.QueryInt("minFreeKilometer", 0),
	}

	// Build the SQL query
	sqlQuery := oc.offerService.buildSQLQuery(params)

	// Execute the query (use your database connection)

	// Further aggregation stuff
	offers := []models.ResponseOffer{}
	// TODO

	// Example aggregations (replace with your database logic)
	priceRanges := []models.PriceRange{{Start: 10000, End: 15000, Count: 4}}
	carTypeCounts := models.CarTypeCounts{Small: 1, Sports: 2, Luxury: 1, Family: 0}
	seatsCount := []models.SeatsCount{{NumberSeats: 5, Count: 4}}
	freeKilometerRange := []models.FreeKilometerRange{{Start: 100, End: 150, Count: 4}}
	vollkaskoCount := models.VollkaskoCount{TrueCount: 3, FalseCount: 1}

	// Prepare the response
	response := fiber.Map{
		"offers":             offers,
		"priceRanges":        priceRanges,
		"carTypeCounts":      carTypeCounts,
		"seatsCount":         seatsCount,
		"freeKilometerRange": freeKilometerRange,
		"vollkaskoCount":     vollkaskoCount,
	}

	return c.JSON(response)

	return c.Status(fiber.StatusNoContent).SendString("TODO")
}*/
func (oc *OfferController) GetOffersHandler(c *fiber.Ctx) error {
	params := models.OfferFilterParams{
		RegionID:              c.QueryInt("regionID"),
		TimeRangeStart:        c.QueryInt("timeRangeStart"),
		TimeRangeEnd:          c.QueryInt("timeRangeEnd"),
		NumberDays:            c.QueryInt("numberDays"),
		SortOrder:             c.Query("sortOrder"),
		Page:                  c.QueryInt("page"),
		PageSize:              c.QueryInt("pageSize"),
		PriceRangeWidth:       c.QueryInt("priceRangeWidth"),
		MinFreeKilometerWidth: c.QueryInt("minFreeKilometerWidth"),
	}

	if minNumberSeats := c.QueryInt("minNumberSeats", 0); minNumberSeats > 0 {
		params.MinNumberSeats = &minNumberSeats
	}

	if minPrice := c.QueryInt("minPrice", 0); minPrice >= 0 {
		params.MinPrice = &minPrice
	}

	if maxPrice := c.QueryInt("maxPrice", -1); maxPrice >= 0 {
		params.MaxPrice = &maxPrice
	}

	if carType := c.Query("carType"); carType != "" {
		params.CarType = &carType
	}

	if onlyVollkasko := c.Query("onlyVollkasko"); onlyVollkasko != "" {
		onlyVollkasko := onlyVollkasko == "true"
		params.OnlyVollkasko = &onlyVollkasko
	}

	if minFreeKilometer := c.QueryInt("minFreeKilometer", -1); minFreeKilometer >= 0 {
		params.MinFreeKilometer = &minFreeKilometer
	}

	response, err := oc.offerService.GetOffers(c, params)
	if err != nil {
		log.Printf("Error fetching offers: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch offers"})
	}

	return c.JSON(response)
}

// CreateOffersHandler verarbeitet die POST-Anfragen
func (oc *OfferController) CreateOffersHandler(c *fiber.Ctx) error {
	//log.Printf("Offers: %v\n", string(c.Body()))
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
