package service

import (
	"context"
	"server/internal/repository"
)

type OfferService struct {
	offerRepository repository.OfferRepository
}

// NewOfferService erstellt einen neuen Service mit dem Repository.
func NewOfferService(repo repository.OfferRepository) *OfferService {
	return &OfferService{offerRepository: repo}
}

// CleanUpOldOffers verwendet das Repository, um alte Angebote zu löschen.
func (s *OfferService) CleanUpOldOffers(ctx context.Context) error {
	return s.offerRepository.DeleteOldOffers(ctx)
}
