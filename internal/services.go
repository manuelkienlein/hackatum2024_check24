package internal

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strconv"
)

// FilterOffers Filter offers based on various parameters
func FilterOffers(dbPool *pgxpool.Pool, c *fiber.Ctx) error {
	// Parse and validate query parameters
	regionID, err := strconv.Atoi(c.Query("regionID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid regionID"})
	}
	timeRangeStart, err := strconv.Atoi(c.Query("timeRangeStart"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid timeRangeStart"})
	}
	timeRangeEnd, err := strconv.Atoi(c.Query("timeRangeEnd"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid timeRangeEnd"})
	}
	numberDays, err := strconv.Atoi(c.Query("numberDays"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid numberDays"})
	}
	sortOrder := c.Query("sortOrder")
	if sortOrder != "price-asc" && sortOrder != "price-desc" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid sortOrder"})
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid page"})
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil || pageSize < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid pageSize"})
	}
	priceRangeWidth, err := strconv.Atoi(c.Query("priceRangeWidth"))
	if err != nil || priceRangeWidth < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid priceRangeWidth"})
	}
	minFreeKilometerWidth, err := strconv.Atoi(c.Query("minFreeKilometerWidth"))
	if err != nil || minFreeKilometerWidth < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid minFreeKilometerWidth"})
	}
	minNumberSeats, _ := strconv.Atoi(c.Query("minNumberSeats"))
	minPrice, _ := strconv.Atoi(c.Query("minPrice"))
	maxPrice, _ := strconv.Atoi(c.Query("maxPrice"))
	carType := c.Query("carType")
	onlyVollkasko, _ := strconv.ParseBool(c.Query("onlyVollkasko"))
	minFreeKilometer, _ := strconv.Atoi(c.Query("minFreeKilometer"))

	// Build SQL query dynamically
	query := `
		SELECT id, data, price, car_type, number_seats, free_kilometers, only_vollkasko
		FROM offers
		WHERE most_specific_region_id = $1
			AND start_date >= $2
			AND end_date <= $3
			AND number_days = $4
			AND price >= $5
			AND price < $6
			AND free_kilometers >= $7
	`
	args := []interface{}{regionID, timeRangeStart, timeRangeEnd, numberDays, minPrice, maxPrice, minFreeKilometer}
	argIdx := len(args)

	// Add dynamic filters
	if minNumberSeats > 0 {
		argIdx++
		query += ` AND number_seats >= $` + strconv.Itoa(argIdx)
		args = append(args, minNumberSeats)
	}
	if carType != "" {
		argIdx++
		query += ` AND car_type = $` + strconv.Itoa(argIdx)
		args = append(args, carType)
	}
	if onlyVollkasko {
		argIdx++
		query += ` AND only_vollkasko = $` + strconv.Itoa(argIdx)
		args = append(args, onlyVollkasko)
	}

	// Add sorting and pagination
	query += ` ORDER BY price ` + sortOrder[6:] + `, id LIMIT $` + strconv.Itoa(argIdx+1) + ` OFFSET $` + strconv.Itoa(argIdx+2)
	args = append(args, pageSize, (page-1)*pageSize)

	// Execute the query
	rows, err := dbPool.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Query execution failed: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch offers"})
	}
	defer rows.Close()

	// Process query results
	offers := make([]ResponseOffer, 0, pageSize)
	priceRangeCounts := make(map[string]int)
	carTypeCounts := CarTypeCounts{}
	seatsCount := make(map[int]int)
	freeKilometerCounts := make(map[string]int)
	vollkaskoCount := VollkaskoCount{}

	for rows.Next() {
		var offer ResponseOffer
		var price, numberSeats, freeKilometers int
		var carType string
		var onlyVollkasko bool

		if err := rows.Scan(&offer.ID, &offer.Data, &price, &carType, &numberSeats, &freeKilometers, &onlyVollkasko); err != nil {
			log.Printf("Row scan failed: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to process offers"})
		}

		// Add offer to list
		offers = append(offers, offer)

		// Aggregate price ranges
		priceRangeKey := fmt.Sprintf("%d-%d", (price/priceRangeWidth)*priceRangeWidth, ((price/priceRangeWidth)+1)*priceRangeWidth)
		priceRangeCounts[priceRangeKey]++

		// Aggregate car type counts
		switch carType {
		case "small":
			carTypeCounts.Small++
		case "sports":
			carTypeCounts.Sports++
		case "luxury":
			carTypeCounts.Luxury++
		case "family":
			carTypeCounts.Family++
		}

		// Aggregate seats count
		seatsCount[numberSeats]++

		// Aggregate free kilometer ranges
		freeKilometerKey := fmt.Sprintf("%d-%d", (freeKilometers/minFreeKilometerWidth)*minFreeKilometerWidth, ((freeKilometers/minFreeKilometerWidth)+1)*minFreeKilometerWidth)
		freeKilometerCounts[freeKilometerKey]++

		// Aggregate vollkasko count
		if onlyVollkasko {
			vollkaskoCount.TrueCount++
		} else {
			vollkaskoCount.FalseCount++
		}
	}

	// Transform aggregated data into required format
	priceRanges := make([]PriceRange, 0, len(priceRangeCounts))
	for key, count := range priceRangeCounts {
		var start, end int
		_, err := fmt.Sscanf(key, "%d-%d", &start, &end)
		if err != nil {
			return err
		}
		priceRanges = append(priceRanges, PriceRange{Start: start, End: end, Count: count})
	}

	seatsCounts := make([]SeatsCount, 0, len(seatsCount))
	for seats, count := range seatsCount {
		seatsCounts = append(seatsCounts, SeatsCount{NumberSeats: seats, Count: count})
	}

	freeKilometerRanges := make([]FreeKilometerRange, 0, len(freeKilometerCounts))
	for key, count := range freeKilometerCounts {
		var start, end int
		_, err := fmt.Sscanf(key, "%d-%d", &start, &end)
		if err != nil {
			return err
		}
		freeKilometerRanges = append(freeKilometerRanges, FreeKilometerRange{Start: start, End: end, Count: count})
	}

	// Return the response
	return c.JSON(fiber.Map{
		"offers":             offers,
		"priceRanges":        priceRanges,
		"carTypeCounts":      carTypeCounts,
		"seatsCount":         seatsCounts,
		"freeKilometerRange": freeKilometerRanges,
		"vollkaskoCount":     vollkaskoCount,
	})
}

// CreateOffers Create a new offer and insert it into the database
func CreateOffers(dbPool *pgxpool.Pool, c *fiber.Ctx) error {
	var request struct {
		Offers []Offer `json:"offers"`
	}

	// Parse the request body
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// Use a single query with connection pooling
	for _, offer := range request.Offers {
		_, err := dbPool.Exec(context.Background(), `
			INSERT INTO offers (id, data, most_specific_region_id, start_date, end_date, number_seats, price, car_type, only_vollkasko, free_kilometers)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			offer.ID, offer.Data, offer.MostSpecificRegionID, offer.StartDate, offer.EndDate, offer.NumberSeats, offer.Price, offer.CarType, offer.OnlyVollkasko, offer.FreeKilometers,
		)
		if err != nil {
			log.Printf("Unable to execute statement: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot execute statement"})
		}
	}

	return c.Status(fiber.StatusOK).SendString("Offers were created successfully")
}

// DeleteOffers Delete outdated offers from the database
func DeleteOffers(dbPool *pgxpool.Pool, c *fiber.Ctx) error {
	// Directly use the connection pool for lightweight operations
	_, err := dbPool.Exec(context.Background(), `
        DELETE FROM offers
        WHERE most_specific_region_id NOT IN (SELECT id FROM static_region_data)
    `)
	if err != nil {
		log.Printf("Unable to delete old offers: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot delete old offers"})
	}

	return c.Status(fiber.StatusOK).SendString("Old offers were cleaned up successfully")
}
