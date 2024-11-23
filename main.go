package main

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type Offer struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type PriceRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
	Count int `json:"count"`
}

type CarTypeCounts struct {
	Small  int `json:"small"`
	Sports int `json:"sports"`
	Luxury int `json:"luxury"`
	Family int `json:"family"`
}

type SeatsCount struct {
	NumberSeats int `json:"numberSeats"`
	Count       int `json:"count"`
}

type FreeKilometerRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
	Count int `json:"count"`
}

type VollkaskoCount struct {
	TrueCount  int `json:"trueCount"`
	FalseCount int `json:"falseCount"`
}

func filterOffers(c *fiber.Ctx, regionID int, timeRangeStart int, timeRangeEnd int, numberDays int, sortOrder string, page int, pageSize int, priceRangeWidth int, minFreeKilometerWidth int, minNumberSeats int, minPrice int, maxPrice int, carType string, onlyVollkasko bool, minFreeKilometer int) error {
	// Placeholder data for response
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

func main() {
	app := fiber.New()

	app.Get("/api/offers", func(c *fiber.Ctx) error {
		regionID, _ := strconv.Atoi(c.Query("regionID"))
		timeRangeStart, _ := strconv.Atoi(c.Query("timeRangeStart"))
		timeRangeEnd, _ := strconv.Atoi(c.Query("timeRangeEnd"))
		numberDays, _ := strconv.Atoi(c.Query("numberDays"))
		sortOrder := c.Query("sortOrder")
		page, _ := strconv.Atoi(c.Query("page"))
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		priceRangeWidth, _ := strconv.Atoi(c.Query("priceRangeWidth"))
		minFreeKilometerWidth, _ := strconv.Atoi(c.Query("minFreeKilometerWidth"))
		minNumberSeats, _ := strconv.Atoi(c.Query("minNumberSeats"))
		minPrice, _ := strconv.Atoi(c.Query("minPrice"))
		maxPrice, _ := strconv.Atoi(c.Query("maxPrice"))
		carType := c.Query("carType")
		onlyVollkasko, _ := strconv.ParseBool(c.Query("onlyVollkasko"))
		minFreeKilometer, _ := strconv.Atoi(c.Query("minFreeKilometer"))

		return filterOffers(c, regionID, timeRangeStart, timeRangeEnd, numberDays, sortOrder, page, pageSize, priceRangeWidth, minFreeKilometerWidth, minNumberSeats, minPrice, maxPrice, carType, onlyVollkasko, minFreeKilometer)
	})

	app.Delete("/api/offers", func(c *fiber.Ctx) error {
		// Logic to clean up the data goes here

		return c.Status(fiber.StatusOK).SendString("Data was cleaned up")
	})

	app.Post("/api/offers", func(c *fiber.Ctx) error {
		// Logic to create the offers goes here

		return c.Status(fiber.StatusOK).SendString("Offers were created")
	})

	app.Listen(":3000")
}
