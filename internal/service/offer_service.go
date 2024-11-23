package service

import (
	"context"
	"server/internal/models"
	"server/internal/repository"
)

type OfferService struct {
	offerRepository repository.OfferRepository
}

// NewOfferService erstellt einen neuen Service mit dem Repository.
func NewOfferService(repo repository.OfferRepository) *OfferService {
	return &OfferService{offerRepository: repo}
}

// CreateOffers erstellt ein neues Offer in der Datenbank
func (s *OfferService) CreateOffers(ctx context.Context, offers []models.Offer) error {
	// Add any business logic here, e.g., validation
	return s.offerRepository.CreateOffers(ctx, offers)
}

// CleanUpOldOffers verwendet das Repository, um alte Angebote zu löschen.
func (s *OfferService) CleanUpOldOffers(ctx context.Context) error {
	return s.offerRepository.DeleteOldOffers(ctx)
}

// Get offers
func (s *OfferService) GetOffers(ctx context.Context, params models.OfferFilterParams) (models.OfferQueryResponse, error) {
	offers, err := s.offerRepository.GetOffers(ctx, params)
	if err != nil {
		return models.OfferQueryResponse{}, err
	}

	// Placeholder: Replace with actual database aggregation logic
	priceRanges := []models.PriceRange{{Start: 10000, End: 15000, Count: 4}}
	carTypeCounts := models.CarTypeCounts{Small: 1, Sports: 2, Luxury: 1, Family: 0}
	seatsCount := []models.SeatsCount{{NumberSeats: 5, Count: 4}}
	freeKilometerRange := []models.FreeKilometerRange{{Start: 100, End: 150, Count: 4}}
	vollkaskoCount := models.VollkaskoCount{TrueCount: 3, FalseCount: 1}

	return models.OfferQueryResponse{
		Offers:             offers,
		PriceRanges:        priceRanges,
		CarTypeCounts:      carTypeCounts,
		SeatsCount:         seatsCount,
		FreeKilometerRange: freeKilometerRange,
		VollkaskoCount:     vollkaskoCount,
	}, nil
}

// Construct a SQL query from filterParams
func (s *OfferService) buildSQLQuery(params models.OfferFilterParams) string {
	// TODO (deprecated)
	return ""
}
