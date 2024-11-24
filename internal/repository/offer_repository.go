package repository

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"server/internal/models"
	"strconv"
	"strings"
)

type OfferRepository interface {
	DeleteOldOffers(ctx context.Context) error
	CreateOffers(ctx context.Context, offers []models.Offer) error
	GetOffers(c *fiber.Ctx, params models.OfferFilterParams) (pgx.Rows, error)
	GetPriceRangeCounts(ctx context.Context, priceRangeWidth int) (map[string]int, error)
	GetCarTypeCounts(ctx context.Context) (models.CarTypeCounts, error)
	GetSeatsCount(ctx context.Context) ([]models.SeatsCount, error)
	GetFreeKilometerCounts(ctx context.Context, minFreeKilometerWidth int) (map[string]int, error)
	GetVollkaskoCount(ctx context.Context) (models.VollkaskoCount, error)
}

type offerRepository struct {
	db *pgxpool.Pool
}

// NewOfferRepository erstellt ein neues Repository mit der Datenbankverbindung.
func NewOfferRepository(db *pgxpool.Pool) OfferRepository {
	return &offerRepository{db: db}
}

// CreateOffers erstellt einen neuen Offer Datensatz in der Datenbank
func (r *offerRepository) CreateOffers(ctx context.Context, offers []models.Offer) error {
	if len(offers) == 0 {
		return nil
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
	  INSERT INTO offers (id, data, most_specific_region_id, start_date, end_date, number_seats, price, car_type, only_vollkasko, free_kilometers)
	  VALUES
	 `)

	args := make([]interface{}, 0, len(offers)*10)
	for i, offer := range offers {
		if i > 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*10+1, i*10+2, i*10+3, i*10+4, i*10+5, i*10+6, i*10+7, i*10+8, i*10+9, i*10+10))
		args = append(args, offer.ID, offer.Data, offer.MostSpecificRegionID, offer.StartDate, offer.EndDate, offer.NumberSeats, offer.Price, offer.CarType, offer.OnlyVollkasko, offer.FreeKilometers)
	}

	query := queryBuilder.String()
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// DeleteOldOffers löscht veraltete Angebote aus der Datenbank.
func (r *offerRepository) DeleteOldOffers(ctx context.Context) error {
	query := `
        DELETE FROM offers
        WHERE end_date < extract(epoch from now())*1000;
    `
	_, err := r.db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to delete old offers: %v\n", err)
		return err
	}
	return nil
}

func (r *offerRepository) GetPriceRangeCounts(ctx context.Context, priceRangeWidth int) (map[string]int, error) {
	query := `
		SELECT (price / $1) * $1 AS range_start, ((price / $1) + 1) * $1 AS range_end, COUNT(*)
		FROM offers
		GROUP BY range_start, range_end
	`
	rows, err := r.db.Query(ctx, query, priceRangeWidth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	priceRangeCounts := make(map[string]int)
	for rows.Next() {
		var rangeStart, rangeEnd, count int
		if err := rows.Scan(&rangeStart, &rangeEnd, &count); err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%d-%d", rangeStart, rangeEnd)
		priceRangeCounts[key] = count
	}
	return priceRangeCounts, nil
}

func (r *offerRepository) GetCarTypeCounts(ctx context.Context) (models.CarTypeCounts, error) {
	query := `
		SELECT car_type, COUNT(*)
		FROM offers
		GROUP BY car_type
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return models.CarTypeCounts{}, err
	}
	defer rows.Close()

	carTypeCounts := models.CarTypeCounts{}
	for rows.Next() {
		var carType string
		var count int
		if err := rows.Scan(&carType, &count); err != nil {
			return carTypeCounts, err
		}
		switch carType {
		case "small":
			carTypeCounts.Small = count
		case "sports":
			carTypeCounts.Sports = count
		case "luxury":
			carTypeCounts.Luxury = count
		case "family":
			carTypeCounts.Family = count
		}
	}
	return carTypeCounts, nil
}

func (r *offerRepository) GetSeatsCount(ctx context.Context) ([]models.SeatsCount, error) {
	query := `
		SELECT number_seats, COUNT(*)
		FROM offers
		GROUP BY number_seats
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seatsCount := []models.SeatsCount{}
	for rows.Next() {
		var numberSeats, count int
		if err := rows.Scan(&numberSeats, &count); err != nil {
			return nil, err
		}
		seatsCount = append(seatsCount, models.SeatsCount{NumberSeats: numberSeats, Count: count})
	}
	return seatsCount, nil
}

func (r *offerRepository) GetFreeKilometerCounts(ctx context.Context, minFreeKilometerWidth int) (map[string]int, error) {
	query := `
		SELECT (free_kilometers / $1) * $1 AS range_start, ((free_kilometers / $1) + 1) * $1 AS range_end, COUNT(*)
		FROM offers
		GROUP BY range_start, range_end
	`
	rows, err := r.db.Query(ctx, query, minFreeKilometerWidth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	freeKilometerCounts := make(map[string]int)
	for rows.Next() {
		var rangeStart, rangeEnd, count int
		if err := rows.Scan(&rangeStart, &rangeEnd, &count); err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%d-%d", rangeStart, rangeEnd)
		freeKilometerCounts[key] = count
	}
	return freeKilometerCounts, nil
}

func (r *offerRepository) GetVollkaskoCount(ctx context.Context) (models.VollkaskoCount, error) {
	query := `
		SELECT only_vollkasko, COUNT(*)
		FROM offers
		GROUP BY only_vollkasko
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return models.VollkaskoCount{}, err
	}
	defer rows.Close()

	vollkaskoCount := models.VollkaskoCount{}
	for rows.Next() {
		var onlyVollkasko bool
		var count int
		if err := rows.Scan(&onlyVollkasko, &count); err != nil {
			return vollkaskoCount, err
		}
		if onlyVollkasko {
			vollkaskoCount.TrueCount = count
		} else {
			vollkaskoCount.FalseCount = count
		}
	}
	return vollkaskoCount, nil
}

func (r *offerRepository) GetOffers(c *fiber.Ctx, params models.OfferFilterParams) (pgx.Rows, error) {
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
				AND o.end_date - start_date >= $4
	`
	args := []interface{}{params.RegionID, params.TimeRangeStart, params.TimeRangeEnd, params.NumberDays * 24 * 3600 * 1000}
	argIdx := len(args)

	// Add dynamic filters
	if params.MinPrice != nil {
		argIdx++
		query += ` AND o.price >= $` + strconv.Itoa(argIdx)
		args = append(args, *params.MinPrice)
	}

	if params.MaxPrice != nil {
		argIdx++
		query += ` AND o.price <= $` + strconv.Itoa(argIdx)
		args = append(args, *params.MaxPrice)
	}

	if params.MinFreeKilometer != nil {
		argIdx++
		query += ` AND o.free_kilometers >= $` + strconv.Itoa(argIdx)
		args = append(args, *params.MinFreeKilometer)
	}

	if params.MinNumberSeats != nil {
		argIdx++
		query += ` AND o.number_seats >= $` + strconv.Itoa(argIdx)
		args = append(args, *params.MinNumberSeats)
	}
	if params.CarType != nil {
		argIdx++
		query += ` AND o.car_type = $` + strconv.Itoa(argIdx)
		args = append(args, *params.CarType)
	}
	if params.OnlyVollkasko != nil {
		argIdx++
		query += ` AND o.only_vollkasko = $` + strconv.Itoa(argIdx)
		args = append(args, *params.OnlyVollkasko)
	}

	// Add sorting and pagination
	query += ` ORDER BY o.price, o.id ` + params.SortOrder[6:] + `, id LIMIT $` + strconv.Itoa(argIdx+1) + ` OFFSET $` + strconv.Itoa(argIdx+2)
	args = append(args, params.PageSize, params.Page*params.PageSize)

	//log.Printf("Query: %v\n", query)
	//log.Printf("SQL query executed: %s, args: %v", query, args)
	//formattedQuery := FormatQuery(query, args)
	//fmt.Println("Formatted Query: ", formattedQuery)

	// Execute the query
	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Query execution failed: %v\n", err)
	}

	return rows, err
}

func FormatQuery(query string, args []interface{}) string {
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		query = strings.Replace(query, placeholder, fmt.Sprintf("'%v'", arg), 1)
	}
	return query
}
