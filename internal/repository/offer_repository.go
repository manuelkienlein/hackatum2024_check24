package repository

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"server/internal/models"
	"strconv"
)

type OfferRepository interface {
	DeleteOldOffers(ctx context.Context) error
	CreateOffers(ctx context.Context, offers []models.Offer) error
	GetOffers(c *fiber.Ctx, params models.OfferFilterParams) (pgx.Rows, error)
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
	for _, offer := range offers {
		_, err := r.db.Exec(ctx, `
			INSERT INTO offers (id, data, most_specific_region_id, start_date, end_date, number_seats, price, car_type, only_vollkasko, free_kilometers)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			offer.ID, offer.Data, offer.MostSpecificRegionID, offer.StartDate, offer.EndDate, offer.NumberSeats, offer.Price, offer.CarType, offer.OnlyVollkasko, offer.FreeKilometers,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteOldOffers löscht veraltete Angebote aus der Datenbank.
func (r *offerRepository) DeleteOldOffers(ctx context.Context) error {
	query := `
        DELETE FROM offers
        WHERE end_date < extract(epoch from now());
    `
	_, err := r.db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to delete old offers: %v\n", err)
		return err
	}
	return nil
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
	args := []interface{}{params.RegionID, params.TimeRangeStart, params.TimeRangeEnd, params.NumberDays}
	argIdx := len(args)

	// Add dynamic filters
	if params.MinPrice != nil {
		argIdx++
		query += ` AND o.price >= $` + strconv.Itoa(argIdx)
		args = append(args, params.MinPrice)
	}

	if params.MaxPrice != nil {
		argIdx++
		query += ` AND o.price <= $` + strconv.Itoa(argIdx)
		args = append(args, params.MaxPrice)
	}

	if params.MinFreeKilometer != nil {
		argIdx++
		query += ` AND o.free_kilometers >= $` + strconv.Itoa(argIdx)
		args = append(args, params.MinFreeKilometer)
	}

	if params.MinNumberSeats != nil {
		argIdx++
		query += ` AND o.number_seats >= $` + strconv.Itoa(argIdx)
		args = append(args, params.MinNumberSeats)
	}
	if params.CarType != nil {
		argIdx++
		query += ` AND o.car_type = $` + strconv.Itoa(argIdx)
		args = append(args, params.CarType)
	}
	if params.OnlyVollkasko != nil {
		argIdx++
		query += ` AND o.only_vollkasko = $` + strconv.Itoa(argIdx)
		args = append(args, params.OnlyVollkasko)
	}

	// Add sorting and pagination
	query += ` ORDER BY o.price ` + params.SortOrder[6:] + `, id LIMIT $` + strconv.Itoa(argIdx+1) + ` OFFSET $` + strconv.Itoa(argIdx+2)
	args = append(args, params.PageSize, params.Page*params.PageSize)

	//log.Printf("Query: %v\n", query)

	// Execute the query
	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Query execution failed: %v\n", err)
	}

	return rows, err
}
