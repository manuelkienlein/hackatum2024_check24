package internal

import "github.com/gofiber/fiber/v2"

func FilterOffers(c *fiber.Ctx, regionID int, timeRangeStart int, timeRangeEnd int, numberDays int, sortOrder string, page int, pageSize int, priceRangeWidth int, minFreeKilometerWidth int, minNumberSeats int, minPrice int, maxPrice int, carType string, onlyVollkasko bool, minFreeKilometer int) error {
	// Placeholder data for response
	// TODO
	offers := []Offer{
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

func CreateOffers(c *fiber.Ctx) error {
	// Placeholder logic to create offers
	// TODO

	return c.Status(fiber.StatusOK).SendString("Offers were created")
}

func DeleteOffers(c *fiber.Ctx) error {
	// Placeholder logic to delete offers
	// TODO

	return c.Status(fiber.StatusOK).SendString("Offers were deleted")
}
