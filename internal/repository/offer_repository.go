package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"server/internal/models"
)

type OfferRepository interface {
	DeleteOldOffers(ctx context.Context) error
	CreateOffers(ctx context.Context, offers []models.Offer) error
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
