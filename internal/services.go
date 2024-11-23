package internal

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// FilterOffers Filter offers based on various parameters
/*func FilterOffers(dbPool *pgxpool.Pool, c *fiber.Ctx) error {
	// print the URL
	log.Println(c.OriginalURL())

	// Parse and validate query parameters
	regionID, err := strconv.Atoi(c.Query("regionID"))
	if err != nil {
		log.Printf("Invalid regionID: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid regionID"})
	}
	//timeRangeStart, err := strconv.Atoi(c.Query("timeRangeStart"))
	if err != nil {
		log.Printf("Invalid timeRangeStart: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid timeRangeStart"})
	}
	//timeRangeEnd, err := strconv.Atoi(c.Query("timeRangeEnd"))
	if err != nil {
		log.Printf("Invalid timeRangeEnd: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid timeRangeEnd"})
	}
	//numberDays, err := strconv.Atoi(c.Query("numberDays"))
	if err != nil {
		log.Printf("Invalid numberDays: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid numberDays"})
	}
	sortOrder := c.Query("sortOrder")
	if sortOrder != "price-asc" && sortOrder != "price-desc" {
		log.Printf("Invalid sortOrder: %v\n", sortOrder)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid sortOrder"})
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 0 {
		log.Printf("Invalid page: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid page"})
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil || pageSize < 1 {
		log.Printf("Invalid pageSize: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid pageSize"})
	}
	priceRangeWidth, err := strconv.Atoi(c.Query("priceRangeWidth"))
	if err != nil || priceRangeWidth < 1 {
		log.Printf("Invalid priceRangeWidth: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid priceRangeWidth"})
	}
	minFreeKilometerWidth, err := strconv.Atoi(c.Query("minFreeKilometerWidth"))
	if err != nil || minFreeKilometerWidth < 1 {
		log.Printf("Invalid minFreeKilometerWidth: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid minFreeKilometerWidth"})
	}
	//minNumberSeats, _ := strconv.Atoi(c.Query("minNumberSeats"))
	//minPrice, _ := strconv.Atoi(c.Query("minPrice"))
	//maxPrice, _ := strconv.Atoi(c.Query("maxPrice"))
	//carType := c.Query("carType")
	//onlyVollkasko, _ := strconv.ParseBool(c.Query("onlyVollkasko"))
	//minFreeKilometer, _ := strconv.Atoi(c.Query("minFreeKilometer"))

	// Build SQL query dynamically
	query := `
		WITH RECURSIVE SubRegions AS (
		-- Base case: Start with the given region ID
		SELECT id, parent_id
		FROM static_region_data
		WHERE id = $1

			UNION ALL

			-- Recursive case: Find children of the current regions
			SELECT sr.id, sr.parent_id
			FROM static_region_data sr
			INNER JOIN SubRegions sbr ON sr.parent_id = sbr.id
		)
		SELECT o.*
		FROM offers o
		JOIN SubRegions sr ON o.most_specific_region_id = sr.id
		WHERE o.start_date >= $2
				AND o.end_date <= $3
				AND o.end_date - start_date = $4
				AND o.price >= $5
				AND o.price < $6
				AND o.free_kilometers >= $7
	`
	args := []interface{}{regionID, timeRangeStart, timeRangeEnd, numberDays, minPrice, maxPrice, minFreeKilometer}
	argIdx := len(args)

	// Add dynamic filters
	if minNumberSeats > 0 {
		argIdx++
		query += ` AND o.number_seats >= $` + strconv.Itoa(argIdx)
		args = append(args, minNumberSeats)
	}
	if carType != "" {
		argIdx++
		query += ` AND o.car_type = $` + strconv.Itoa(argIdx)
		args = append(args, carType)
	}
	if onlyVollkasko {
		argIdx++
		query += ` AND o.only_vollkasko = $` + strconv.Itoa(argIdx)
		args = append(args, onlyVollkasko)
	}

	// Add sorting and pagination
	query += ` ORDER BY o.price ` + sortOrder[6:] + `, id LIMIT $` + strconv.Itoa(argIdx+1) + ` OFFSET $` + strconv.Itoa(argIdx+2)
	args = append(args, pageSize, (page-1)*pageSize)

	log.Printf("Query: %v\n", query)

	// Execute the query
	rows, err := dbPool.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Query execution failed: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch offers"})
	}
	defer rows.Close()

	// Process query results
	offers := make([]map[string]string, 0, pageSize)
	priceRangeCounts := make(map[string]int)
	carTypeCounts := map[string]int{"small": 0, "sports": 0, "luxury": 0, "family": 0}
	seatsCount := make(map[int]int)
	freeKilometerCounts := make(map[string]int)
	vollkaskoCount := map[string]int{"trueCount": 0, "falseCount": 0}

	for rows.Next() {
		var id, data, carType string
		var regionId, startDate, endDate, price, numberSeats, freeKilometers int
		var onlyVollkasko bool

		if err := rows.Scan(&id, &data, &regionId, &startDate, &endDate, &numberSeats, &price, &carType, &onlyVollkasko, &freeKilometers); err != nil {
			log.Printf("Row scan failed: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to process offers"})
		}

		// Add offer to list
		offers = append(offers, map[string]string{"ID": id, "data": data})

		// Aggregate price ranges
		priceRangeKey := fmt.Sprintf("%d-%d", (price/priceRangeWidth)*priceRangeWidth, ((price/priceRangeWidth)+1)*priceRangeWidth)
		priceRangeCounts[priceRangeKey]++

		// Aggregate car type counts
		carTypeCounts[carType]++

		// Aggregate seats count
		seatsCount[numberSeats]++

		// Aggregate free kilometer ranges
		freeKilometerKey := fmt.Sprintf("%d-%d", (freeKilometers/minFreeKilometerWidth)*minFreeKilometerWidth, ((freeKilometers/minFreeKilometerWidth)+1)*minFreeKilometerWidth)
		freeKilometerCounts[freeKilometerKey]++

		// Aggregate vollkasko count
		if onlyVollkasko {
			vollkaskoCount["trueCount"]++
		} else {
			vollkaskoCount["falseCount"]++
		}
	}

	// Transform aggregated data into required format
	priceRanges := make([]map[string]int, 0, len(priceRangeCounts))
	for key, count := range priceRangeCounts {
		var start, end int
		_, err := fmt.Sscanf(key, "%d-%d", &start, &end)
		if err != nil {
			return err
		}
		priceRanges = append(priceRanges, map[string]int{"start": start, "end": end, "count": count})
	}

	seatsCounts := make([]map[string]int, 0, len(seatsCount))
	for seats, count := range seatsCount {
		seatsCounts = append(seatsCounts, map[string]int{"numberSeats": seats, "count": count})
	}

	freeKilometerRanges := make([]map[string]int, 0, len(freeKilometerCounts))
	for key, count := range freeKilometerCounts {
		var start, end int
		_, err := fmt.Sscanf(key, "%d-%d", &start, &end)
		if err != nil {
			return err
		}
		freeKilometerRanges = append(freeKilometerRanges, map[string]int{"start": start, "end": end, "count": count})
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
}*/

// CreateOffers Create a new offer and insert it into the database
func CreateOffers(dbPool *pgxpool.Pool, c *fiber.Ctx) error {
	// print out the request body
	log.Println(string(c.Body()))

	var request struct {
		Offers []Offer `json:"offers"`
	}

	// Parse the request body
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Unable to parse JSON: %v\n", err)
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
