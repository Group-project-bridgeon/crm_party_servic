package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/repository"
)

type OrganizationService struct {
	partyRepository        repository.PartyRepository
	organizationRepository repository.OrganizationRepository
}

func NewOrganizationService(
	partyRepository repository.PartyRepository,
	organizationRepository repository.OrganizationRepository,
) *OrganizationService {
	return &OrganizationService{
		partyRepository:        partyRepository,
		organizationRepository: organizationRepository,
	}
}

func (s *OrganizationService) CreateOrganization(
	ctx context.Context,
	input domain.CreateOrganizationInput,
) (db.Organization, error) {
	// 1. Parse party ID.
	var partyID pgtype.UUID

	if err := partyID.Scan(input.PartyID); err != nil {
		return db.Organization{}, ErrInvalidPartyID
	}

	// 2. Check that the party exists.
	party, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Organization{}, ErrPartyNotFound
		}

		return db.Organization{}, err
	}

	// 3. Make sure the party is an ORGANIZATION.
	if party.PartyType != string(domain.PartyTypeOrganization) {
		return db.Organization{}, ErrPartyNotOrganization
	}

	// 4. Make sure an organization doesn't already exist.
	_, err = s.organizationRepository.GetOrganizationByPartyID(ctx, partyID)
	if err == nil {
		return db.Organization{}, ErrOrganizationAlreadyExists
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return db.Organization{}, err
	}

	// 5. Validate employee count.
	if input.EmployeeCount != nil {
		if *input.EmployeeCount < 0 {
			return db.Organization{}, ErrInvalidEmployeeCount
		}

		if int64(*input.EmployeeCount) > math.MaxInt32 {
			return db.Organization{}, ErrInvalidEmployeeCount
		}
	}

	// 6. Convert founded date.
	var foundedDate pgtype.Date

	if input.FoundedDate != nil && *input.FoundedDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *input.FoundedDate)
		if err != nil {
			return db.Organization{}, ErrInvalidFoundedDate
		}

		foundedDate = pgtype.Date{
			Time:  parsedDate,
			Valid: true,
		}
	}

	// 7. Convert employee count from int to int32.
	var employeeCount *int32

	if input.EmployeeCount != nil {
		value := int32(*input.EmployeeCount)
		employeeCount = &value
	}

	// 8. Build database parameters.
	params := db.CreateOrganizationParams{
		PartyID:            partyID,
		RegistrationNumber: input.RegistrationNumber,
		Industry:           input.Industry,
		OrganizationType:   input.OrganizationType,
		FoundedDate:        foundedDate,
		EmployeeCount:      employeeCount,
		Description:        input.Description,
	}

	// 9. Create organization.
	return s.organizationRepository.CreateOrganization(ctx, params)
}


func (s *OrganizationService) GetOrganizationByPartyID(
    ctx context.Context,
    partyID pgtype.UUID,
) (db.Organization, error) {
    organization, err := s.organizationRepository.GetOrganizationByPartyID(ctx, partyID)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return db.Organization{}, ErrOrganizationNotFound
        }

        return db.Organization{}, err
    }

    return organization, nil
}

func (s *OrganizationService) UpdateOrganization(
	ctx context.Context,
	input domain.UpdateOrganizationInput,
) (db.Organization, error) {
	// 1. Parse party ID.
	var partyID pgtype.UUID
	if err := partyID.Scan(input.PartyID); err != nil {
		return db.Organization{}, ErrInvalidPartyID
	}

	// 2. Check that the party exists.
	party, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Organization{}, ErrPartyNotFound
		}
		return db.Organization{}, err
	}

	// 3. Make sure the party is an ORGANIZATION.
	if party.PartyType != string(domain.PartyTypeOrganization) {
		return db.Organization{}, ErrPartyNotOrganization
	}

	// 4. Make sure an organization profile exists.
	_, err = s.organizationRepository.GetOrganizationByPartyID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Organization{}, ErrOrganizationNotFound
		}
		return db.Organization{}, err
	}

	// 5. Validate employee count.
	if input.EmployeeCount != nil {
		if *input.EmployeeCount < 0 || int64(*input.EmployeeCount) > math.MaxInt32 {
			return db.Organization{}, ErrInvalidEmployeeCount
		}
	}

	// 6. Convert founded date.
	var foundedDate pgtype.Date
	if input.FoundedDate != nil && *input.FoundedDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *input.FoundedDate)
		if err != nil {
			return db.Organization{}, ErrInvalidFoundedDate
		}
		foundedDate = pgtype.Date{
			Time:  parsedDate,
			Valid: true,
		}
	}

	// 7. Convert employee count from int to int32.
	var employeeCount *int32
	if input.EmployeeCount != nil {
		value := int32(*input.EmployeeCount)
		employeeCount = &value
	}

	params := db.UpdateOrganizationParams{
		PartyID:            partyID,
		RegistrationNumber: input.RegistrationNumber,
		Industry:           input.Industry,
		OrganizationType:   input.OrganizationType,
		FoundedDate:        foundedDate,
		EmployeeCount:      employeeCount,
		Description:        input.Description,
	}

	return s.organizationRepository.UpdateOrganization(ctx, params)
}
