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
	priceRangeCounts := make(map[string]int)
	carTypeCounts := models.CarTypeCounts{Small: 0, Sports: 0, Luxury: 0, Family: 0}
	seatsCount := make([]models.SeatsCount, 0)
	freeKilometerCounts := make(map[string]int)
	vollkaskoCount := models.VollkaskoCount{TrueCount: 0, FalseCount: 0}

	for rows.Next() {
		var id, data, carType string
		var regionId, startDate, endDate, price, numberSeats, freeKilometers int
		var onlyVollkasko bool

		if err := rows.Scan(&id, &data, &regionId, &startDate, &endDate, &numberSeats, &price, &carType, &onlyVollkasko, &freeKilometers); err != nil {
			log.Printf("Row scan failed: %v\n", err)
			return models.OfferQueryResponse{}, err
		}

		// Check aggregate filters
		minNumberSeatsFlag := true
		minPriceFlag := true
		maxPriceFlag := true
		carTypeFlag := params.CarType == nil || carType == *params.CarType
		onlyVollkaskoFlag := true
		minFreeKilometerFlag := true

		// if all aggregate filters are satisfied (and not nil), add the offer to the response
		if minNumberSeatsFlag && minPriceFlag && maxPriceFlag && carTypeFlag && onlyVollkaskoFlag && minFreeKilometerFlag {
			offers = append(offers, models.ResponseOffer{ID: id, Data: data})
		}

		// Aggregate price ranges
		if minNumberSeatsFlag && carTypeFlag && onlyVollkaskoFlag && minFreeKilometerFlag {
			priceRangeKey := fmt.Sprintf("%d-%d", (price/params.PriceRangeWidth)*params.PriceRangeWidth, ((price/params.PriceRangeWidth)+1)*params.PriceRangeWidth)
			priceRangeCounts[priceRangeKey]++
		}

		// Aggregate car type counts
		if minNumberSeatsFlag && minPriceFlag && maxPriceFlag && onlyVollkaskoFlag && minFreeKilometerFlag {
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
		}

		// Aggregate seats count
		if minPriceFlag && maxPriceFlag && carTypeFlag && onlyVollkaskoFlag && minFreeKilometerFlag {
			found := false
			for i, sc := range seatsCount {
				if sc.NumberSeats == numberSeats {
					seatsCount[i].Count++
					found = true
					break
				}
			}
			if !found {
				seatsCount = append(seatsCount, models.SeatsCount{NumberSeats: numberSeats, Count: 1})
			}
		}

		// Aggregate free kilometer ranges
		if minNumberSeatsFlag && minPriceFlag && maxPriceFlag && carTypeFlag && onlyVollkaskoFlag {
			freeKilometerKey := fmt.Sprintf("%d-%d", (freeKilometers/params.MinFreeKilometerWidth)*params.MinFreeKilometerWidth, ((freeKilometers/params.MinFreeKilometerWidth)+1)*params.MinFreeKilometerWidth)
			freeKilometerCounts[freeKilometerKey]++
		}

		// Aggregate vollkasko count
		if minNumberSeatsFlag && minPriceFlag && maxPriceFlag && carTypeFlag && minFreeKilometerFlag {
			if onlyVollkasko {
				vollkaskoCount.TrueCount++
			} else {
				vollkaskoCount.FalseCount++
			}
		}
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
