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
