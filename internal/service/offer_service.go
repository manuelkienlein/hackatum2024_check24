package service

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
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
func (s *OfferService) GetOffers(c *fiber.Ctx, params models.OfferFilterParams) (models.OfferQueryResponse, error) {
	rows, err := s.offerRepository.GetOffers(c, params)
	defer rows.Close()
	if err != nil {
		return models.OfferQueryResponse{}, err
	}

	// Process query results
	offers := make([]models.ResponseOffer, 0, params.PageSize)
	for rows.Next() {
		var id, data string
		var regionId, startDate, endDate, price, numberSeats, freeKilometers int
		var carType string
		var onlyVollkasko bool

		if err := rows.Scan(&id, &data, &regionId, &startDate, &endDate, &numberSeats, &price, &carType, &onlyVollkasko, &freeKilometers); err != nil {
			log.Printf("Row scan failed: %v\n", err)
			return models.OfferQueryResponse{}, err
		}

		// Add offer to list
		offers = append(offers, models.ResponseOffer{ID: id, Data: data})
	}

	// Fetch counts for the whole database
	ctx := context.Background()
	priceRangeCounts, err := s.offerRepository.GetPriceRangeCounts(ctx, params.PriceRangeWidth)
	if err != nil {
		return models.OfferQueryResponse{}, err
	}

	carTypeCounts, err := s.offerRepository.GetCarTypeCounts(ctx)
	if err != nil {
		return models.OfferQueryResponse{}, err
	}

	seatsCount, err := s.offerRepository.GetSeatsCount(ctx)
	if err != nil {
		return models.OfferQueryResponse{}, err
	}

	freeKilometerCounts, err := s.offerRepository.GetFreeKilometerCounts(ctx, params.MinFreeKilometerWidth)
	if err != nil {
		return models.OfferQueryResponse{}, err
	}

	vollkaskoCount, err := s.offerRepository.GetVollkaskoCount(ctx)
	if err != nil {
		return models.OfferQueryResponse{}, err
	}

	// Transform aggregated data into required format
	priceRanges := make([]models.PriceRange, 0, len(priceRangeCounts))
	for key, count := range priceRangeCounts {
		var start, end int
		_, err := fmt.Sscanf(key, "%d-%d", &start, &end)
		if err != nil {
			return models.OfferQueryResponse{}, err
		}
		priceRanges = append(priceRanges, models.PriceRange{Start: start, End: end, Count: count})
	}

	freeKilometerRanges := make([]models.FreeKilometerRange, 0, len(freeKilometerCounts))
	for key, count := range freeKilometerCounts {
		var start, end int
		_, err := fmt.Sscanf(key, "%d-%d", &start, &end)
		if err != nil {
			return models.OfferQueryResponse{}, err
		}
		freeKilometerRanges = append(freeKilometerRanges, models.FreeKilometerRange{Start: start, End: end, Count: count})
	}

	// Return the response
	return models.OfferQueryResponse{
		Offers:             offers,
		PriceRanges:        priceRanges,
		CarTypeCounts:      carTypeCounts,
		SeatsCount:         seatsCount,
		FreeKilometerRange: freeKilometerRanges,
		VollkaskoCount:     vollkaskoCount,
	}, nil
}

// Construct a SQL query from filterParams
func (s *OfferService) buildSQLQuery(params models.OfferFilterParams) string {
	// TODO (deprecated)
	return ""
}
