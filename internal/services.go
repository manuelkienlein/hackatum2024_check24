package internal

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func FilterOffers(c *fiber.Ctx, regionID int, timeRangeStart int, timeRangeEnd int, numberDays int, sortOrder string, page int, pageSize int, priceRangeWidth int, minFreeKilometerWidth int, minNumberSeats int, minPrice int, maxPrice int, carType string, onlyVollkasko bool, minFreeKilometer int) error {
	// Placeholder data for response
	// TODO
	offers := []ResponseOffer{
		{ID: "offer1", Data: "data1"},
		{ID: "offer2", Data: "data2"},
	}

	priceRanges := []PriceRange{
		{Start: 0, End: 1000, Count: 5},
		{Start: 1000, End: 2000, Count: 3},
	}

	carTypeCounts := CarTypeCounts{
		Small:  2,
		Sports: 1,
		Luxury: 1,
		Family: 1,
	}

	seatsCount := []SeatsCount{
		{NumberSeats: 4, Count: 3},
		{NumberSeats: 5, Count: 2},
	}

	freeKilometerRange := []FreeKilometerRange{
		{Start: 0, End: 100, Count: 4},
		{Start: 100, End: 200, Count: 2},
	}

	vollkaskoCount := VollkaskoCount{
		TrueCount:  3,
		FalseCount: 2,
	}

	return c.JSON(fiber.Map{
		"offers":             offers,
		"priceRanges":        priceRanges,
		"carTypeCounts":      carTypeCounts,
		"seatsCount":         seatsCount,
		"freeKilometerRange": freeKilometerRange,
		"vollkaskoCount":     vollkaskoCount,
	})
}

func CreateOffers(dbPool *pgxpool.Pool, c *fiber.Ctx) error {
	var offer Offer

	// Parse the request body
	if err := c.BodyParser(&offer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// Use the connection pool to get a connection
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot acquire database connection"})
	}
	defer conn.Release()

	// Insert the offer into the database
	_, err = conn.Exec(context.Background(), `
  INSERT INTO offers (region_id, time_range_start, time_range_end, number_days, sort_order, page, page_size, price_range_width, min_free_kilometer_width, min_number_seats, min_price, max_price, car_type, only_vollkasko, min_free_kilometer)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		offer.RegionID, offer.TimeRangeStart, offer.TimeRangeEnd, offer.NumberDays, offer.SortOrder, offer.Page, offer.PageSize, offer.PriceRangeWidth, offer.MinFreeKilometerWidth, offer.MinNumberSeats, offer.MinPrice, offer.MaxPrice, offer.CarType, offer.OnlyVollkasko, offer.MinFreeKilometer)
	if err != nil {
		log.Printf("Unable to insert offer: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot insert offer"})
	}

	return c.Status(fiber.StatusOK).SendString("Offers were created")
}

func DeleteOffers(c *fiber.Ctx) error {
	// Placeholder logic to delete offers
	// TODO

	return c.Status(fiber.StatusOK).SendString("Offers were deleted")
}
