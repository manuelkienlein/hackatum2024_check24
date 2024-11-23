package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"server/internal/models"
)

type OfferRepository interface {
	DeleteOldOffers(ctx context.Context) error
	CreateOffers(ctx context.Context, offers []models.Offer) error
	GetOffers(ctx context.Context, params models.OfferFilterParams) ([]models.ResponseOffer, error)
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

func (r *offerRepository) GetOffers(ctx context.Context, params models.OfferFilterParams) ([]models.ResponseOffer, error) {
	query := buildSQLQuery(params)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var offers []models.ResponseOffer
	for rows.Next() {
		var offer models.ResponseOffer
		if err := rows.Scan(&offer.ID, &offer.Data); err != nil {
			return nil, fmt.Errorf("failed to parse query results: %w", err)
		}
		offers = append(offers, offer)
	}

	return offers, nil
}

// Helper function to build the SQL query
func buildSQLQuery(params models.OfferFilterParams) string {
	query := "SELECT ID, data FROM offers WHERE 1=1"

	query += fmt.Sprintf(" AND mostSpecificRegionID = %d", params.RegionID)
	query += fmt.Sprintf(" AND startDate >= %d", params.TimeRangeStart)
	query += fmt.Sprintf(" AND endDate <= %d", params.TimeRangeEnd)
	query += fmt.Sprintf(" AND julianday(endDate, 'unixepoch') - julianday(startDate, 'unixepoch') >= %d", params.NumberDays)

	if params.MinPrice > 0 {
		query += fmt.Sprintf(" AND price >= %d", params.MinPrice)
	}
	if params.MaxPrice > 0 {
		query += fmt.Sprintf(" AND price < %d", params.MaxPrice)
	}
	if params.MinNumberSeats > 0 {
		query += fmt.Sprintf(" AND numberSeats >= %d", params.MinNumberSeats)
	}
	if params.MinFreeKilometer > 0 {
		query += fmt.Sprintf(" AND freeKilometers >= %d", params.MinFreeKilometer)
	}
	if params.CarType != "" {
		query += fmt.Sprintf(" AND carType = '%s'", params.CarType)
	}
	if params.OnlyVollkasko {
		query += " AND hasVollkasko = 1"
	}

	if params.SortOrder == "price-asc" {
		query += " ORDER BY price ASC, ID ASC"
	} else if params.SortOrder == "price-desc" {
		query += " ORDER BY price DESC, ID ASC"
	}

	offset := (params.Page - 1) * params.PageSize
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.PageSize, offset)

	return query
}
