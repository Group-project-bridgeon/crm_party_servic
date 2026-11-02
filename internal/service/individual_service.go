package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/repository"
)

type IndividualService struct {
	partyRepository      repository.PartyRepository
	individualRepository repository.IndividualRepository
}

func NewIndividualService(
	partyRepository repository.PartyRepository,
	individualRepository repository.IndividualRepository,
) *IndividualService {
	return &IndividualService{
		partyRepository:      partyRepository,
		individualRepository: individualRepository,
	}
}

func (s *IndividualService) CreateIndividual(
	ctx context.Context,
	input domain.CreateIndividualInput,
) (db.Individual, error) {
	// 1. Parse party ID.
	var partyID pgtype.UUID
	if err := partyID.Scan(input.PartyID); err != nil {
		return db.Individual{}, ErrInvalidPartyID
	}

	// 2. Validate first name.
	input.FirstName = strings.TrimSpace(input.FirstName)
	if input.FirstName == "" {
		return db.Individual{}, ErrEmptyFirstName
	}

	// 3. Check party exists.
	party, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Individual{}, ErrPartyNotFound
		}
		return db.Individual{}, err
	}

	// 4. Ensure party is an INDIVIDUAL.
	if party.PartyType != string(domain.PartyTypeIndividual) {
		return db.Individual{}, ErrPartyNotIndividual
	}

	// 5. Ensure an individual doesn't already exist.
	_, err = s.individualRepository.GetIndividualByPartyID(ctx, partyID)
	if err == nil {
		return db.Individual{}, ErrIndividualAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return db.Individual{}, err
	}

	// 6. Convert date of birth.
	var dob pgtype.Date
	if input.DateOfBirth != nil && *input.DateOfBirth != "" {
		parsedDate, err := time.Parse("2006-01-02", *input.DateOfBirth)
		if err != nil {
			return db.Individual{}, ErrInvalidDateOfBirth
		}
		dob = pgtype.Date{
			Time:  parsedDate,
			Valid: true,
		}
	}

	params := db.CreateIndividualParams{
		PartyID:     partyID,
		FirstName:   input.FirstName,
		MiddleName:  input.MiddleName,
		LastName:    input.LastName,
		DateOfBirth: dob,
		Gender:      input.Gender,
		JobTitle:    input.JobTitle,
		NationalID:  input.NationalID,
	}

	return s.individualRepository.CreateIndividual(ctx, params)
}

func (s *IndividualService) GetIndividualByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) (db.Individual, error) {
	individual, err := s.individualRepository.GetIndividualByPartyID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Individual{}, ErrIndividualNotFound
		}
		return db.Individual{}, err
	}

	return individual, nil
}

func (s *IndividualService) UpdateIndividual(
	ctx context.Context,
	input domain.UpdateIndividualInput,
) (db.Individual, error) {
	// 1. Parse party ID.
	var partyID pgtype.UUID
	if err := partyID.Scan(input.PartyID); err != nil {
		return db.Individual{}, ErrInvalidPartyID
	}

	// 2. Validate first name.
	input.FirstName = strings.TrimSpace(input.FirstName)
	if input.FirstName == "" {
		return db.Individual{}, ErrEmptyFirstName
	}

	// 3. Check party exists.
	party, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Individual{}, ErrPartyNotFound
		}
		return db.Individual{}, err
	}

	// 4. Ensure party is an INDIVIDUAL.
	if party.PartyType != string(domain.PartyTypeIndividual) {
		return db.Individual{}, ErrPartyNotIndividual
	}

	// 5. Ensure individual profile exists.
	_, err = s.individualRepository.GetIndividualByPartyID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Individual{}, ErrIndividualNotFound
		}
		return db.Individual{}, err
	}

	// 6. Convert date of birth.
	var dob pgtype.Date
	if input.DateOfBirth != nil && *input.DateOfBirth != "" {
		parsedDate, err := time.Parse("2006-01-02", *input.DateOfBirth)
		if err != nil {
			return db.Individual{}, ErrInvalidDateOfBirth
		}
		dob = pgtype.Date{
			Time:  parsedDate,
			Valid: true,
		}
	}

	params := db.UpdateIndividualParams{
		PartyID:     partyID,
		FirstName:   input.FirstName,
		MiddleName:  input.MiddleName,
		LastName:    input.LastName,
		DateOfBirth: dob,
		Gender:      input.Gender,
		JobTitle:    input.JobTitle,
		NationalID:  input.NationalID,
	}

	return s.individualRepository.UpdateIndividual(ctx, params)
}
