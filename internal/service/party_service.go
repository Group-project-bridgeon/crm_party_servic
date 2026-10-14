package service

import (
	"context"
	"strings"

	"errors"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
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
	party, err := s.repository.GetPartyByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Party{}, ErrPartyNotFound
		}

		return db.Party{}, err
	}

	return party, nil
}

func (s *PartyService) ListParties(
	ctx context.Context,
	params db.ListPartiesParams,
) ([]db.Party, error) {
	if params.Offset < 0 {
		return nil, ErrInvalidPagination
	}

	if params.Limit <= 0 {
		params.Limit = 20
	}

	if params.Limit > 100 {
		params.Limit = 100
	}

	if params.Search != nil {
		search := strings.TrimSpace(*params.Search)

		if search == "" {
			params.Search = nil
		} else {
			params.Search = &search
		}
	}

	return s.repository.ListParties(ctx, params)
}

func (s *PartyService) CreateParty(
	ctx context.Context,
	input domain.CreatePartyInput,
) (db.Party, error) {
	if !input.PartyType.IsValid() {
		return db.Party{}, ErrInvalidPartyType
	}

	if strings.TrimSpace(input.DisplayName) == "" {
		return db.Party{}, ErrEmptyDisplayName
	}

	input.DisplayName = strings.TrimSpace(input.DisplayName)

	params := db.CreatePartyParams{
		PartyType:   string(input.PartyType),
		DisplayName: input.DisplayName,
		LegalName:   input.LegalName,
		Email:       input.Email,
		Phone:       input.Phone,
		Website:     input.Website,
		TaxNumber:   input.TaxNumber,
		Notes:       input.Notes,
	}

	return s.repository.CreateParty(ctx, params)
}
