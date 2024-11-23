package internal

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func RegisterRoutes(app *fiber.App) {

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

		return FilterOffers(c, regionID, timeRangeStart, timeRangeEnd, numberDays, sortOrder, page, pageSize, priceRangeWidth, minFreeKilometerWidth, minNumberSeats, minPrice, maxPrice, carType, onlyVollkasko, minFreeKilometer)
	})

	app.Delete("/api/offers", func(c *fiber.Ctx) error {
		// Logic to clean up the data goes here

		return DeleteOffers(c)
	})

	app.Post("/api/offers", func(c *fiber.Ctx) error {
		// Logic to create the offers goes here

		return CreateOffers(c)
	})

}
