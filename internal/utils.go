package internal

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func ParseOfferQueryParams(c *fiber.Ctx) (OfferFilterParams, error) {
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

	return OfferFilterParams{
		RegionID:              regionID,
		TimeRangeStart:        timeRangeStart,
		TimeRangeEnd:          timeRangeEnd,
		NumberDays:            numberDays,
		SortOrder:             sortOrder,
		Page:                  page,
		PageSize:              pageSize,
		PriceRangeWidth:       priceRangeWidth,
		MinFreeKilometerWidth: minFreeKilometerWidth,
		MinNumberSeats:        minNumberSeats,
		MinPrice:              minPrice,
		MaxPrice:              maxPrice,
		CarType:               carType,
		OnlyVollkasko:         onlyVollkasko,
		MinFreeKilometer:      minFreeKilometer,
	}, nil
}
