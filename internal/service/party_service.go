package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/repository"
)

type PartyService struct {
	repository repository.PartyRepository
}

func NewPartyService(repository repository.PartyRepository) *PartyService {
	return &PartyService{
		repository: repository,
	}
}

func (s *PartyService) GetPartyByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.Party, error) {
	return s.repository.GetPartyByID(ctx, id)
}

func (s *PartyService) ListParties(
	ctx context.Context,
	params db.ListPartiesParams,
) ([]db.Party, error) {
	return s.repository.ListParties(ctx, params)
}

func (s *PartyService) CreateParty(
	ctx context.Context,
	params db.CreatePartyParams,
) (db.Party, error) {
	return s.repository.CreateParty(ctx, params)
}
